"""
An executable which is run for each connection.

This treats stdin/stdout as a connection to a client.
"""

from argparse import ArgumentParser
import io
import json
import sys

import proto
import gym
from gym import wrappers
import universe_plugin

def main():
    """
    Executable entry-point.
    """
    parser = ArgumentParser()
    parser.add_argument('--addr', action='store', type=str, dest='addr')
    parser.add_argument('--fd', action='store', type=int, dest='fd')
    parser.add_argument('--universe', action='store_true', dest='universe')
    parser.add_argument('--setup', action='store', type=str, dest='setup_code')
    options = parser.parse_args()

    # pylint: disable=W0122
    exec(options.setup_code)

    in_file = io.open(options.fd, 'rb', buffering=0)
    out_file = io.open(options.fd, 'wb', buffering=0)
    handle(io.BufferedRWPair(in_file, out_file), options)

def handle(sock_file, info):
    """
    Handle a connection from a client.
    """
    try:
        uni = universe_plugin.Universe(info.universe)
        env = handshake(sock_file)
        try:
            loop(sock_file, uni, env)
        finally:
            if not env is None:
                env.close()
    except proto.ProtoException as exc:
        log('%s gave error: %s' % (info.addr, str(exc)))

def handshake(sock):
    """
    Perform the initial handshake and return the resulting
    Gym environment.
    """
    flags = proto.read_flags(sock)
    if flags != 0:
        raise proto.ProtoException('unsupported flags: ' + str(flags))
    env_name = proto.read_field_str(sock)

    # Special no-environment mode.
    if env_name == '':
        proto.write_field_str(sock, '')
        sock.flush()
        return None

    try:
        env = gym.make(env_name)
        proto.write_field_str(sock, '')
        sock.flush()
        return env
    except gym.error.Error as gym_exc:
        proto.write_field_str(sock, str(gym_exc))
        sock.flush()
        raise gym_exc

def loop(sock, uni, env):
    """
    Handle commands from the client as they come in and
    apply them to the given Gym environment.
    """
    while True:
        pack_type = proto.read_packet_type(sock)
        if pack_type == 'reset':
            handle_reset(sock, env)
        elif pack_type == 'step':
            handle_step(sock, env)
        elif pack_type == 'get_space':
            handle_get_space(sock, env)
        elif pack_type == 'sample_action':
            handle_sample_action(sock, env)
        elif pack_type == 'monitor':
            env = handle_monitor(sock, env)
        elif pack_type == 'render':
            handle_render(env)
        elif pack_type == 'upload':
            handle_upload(sock)
        elif pack_type == 'universe_configure':
            env = handle_universe_configure(sock, uni, env)
        elif pack_type == 'universe_wrap':
            env = handle_universe_wrap(sock, uni, env)

def handle_reset(sock, env):
    """
    Reset the environment and send the result.
    """
    proto.write_obs(sock, env, env.reset())
    sock.flush()

def handle_step(sock, env):
    """
    Step the environment and send the result.
    """
    action = proto.read_action(sock, env)
    obs, rew, done, info = env.step(action)
    proto.write_obs(sock, env, obs)
    proto.write_reward(sock, rew)
    proto.write_bool(sock, done)
    try:
        dumped_info = json.dumps(info)
    except TypeError:
        dumped_info = '{}'
    proto.write_field_str(sock, dumped_info)
    sock.flush()

def handle_get_space(sock, env):
    """
    Get information about the action or observation space.
    """
    space_id = proto.read_space_id(sock)
    if space_id == 'action':
        proto.write_space(sock, env.action_space)
    elif space_id == 'observation':
        proto.write_space(sock, env.observation_space)
    sock.flush()

def handle_sample_action(sock, env):
    """
    Generate and send a random action.
    """
    action = env.action_space.sample()
    proto.write_action(sock, env, action)
    sock.flush()

def handle_monitor(sock, env):
    """
    Start a monitor and return the new environment.
    """
    resume = proto.read_bool(sock)
    force = proto.read_bool(sock)
    video = proto.read_bool(sock)
    dir_path = proto.read_field_str(sock)
    try:
        vid_call = None
        if not video:
            vid_call = lambda count: False
        res = wrappers.Monitor(env, dir_path, resume=resume, force=force,
                               video_callable=vid_call)
        proto.write_field_str(sock, '')
        sock.flush()
        return res
    except gym.error.Error as exc:
        proto.write_field_str(sock, str(exc))
        sock.flush()
        return env

def handle_render(env):
    """
    Render the environment.
    """
    env.render()

def handle_upload(sock):
    """
    Upload a monitor to the Gym website.
    """
    dir_path = proto.read_field_str(sock)
    api_key = proto.read_field_str(sock)
    alg_id = proto.read_field_str(sock)
    if alg_id == '':
        alg_id = None
    try:
        gym.upload(dir_path, api_key=api_key, algorithm_id=alg_id)
        proto.write_field_str(sock, '')
    except gym.error.Error as exc:
        proto.write_field_str(sock, str(exc))
    sock.flush()

def handle_universe_configure(sock, uni, env):
    """
    Configure a Universe environment.
    """
    config_json = proto.read_field_str(sock)
    try:
        env = uni.configure(env, json.loads(config_json))
        proto.write_field_str(sock, '')
    except universe_plugin.UniverseException as exc:
        proto.write_field_str(sock, str(exc))
    sock.flush()
    return env

def handle_universe_wrap(sock, uni, env):
    """
    Wrap a Universe environment.
    """
    wrapper_name = proto.read_field_str(sock)
    config_json = proto.read_field_str(sock)
    try:
        env = uni.wrap(env, wrapper_name, json.loads(config_json))
        proto.write_field_str(sock, '')
    except universe_plugin.UniverseException as exc:
        proto.write_field_str(sock, str(exc))
    sock.flush()
    return env

def log(msg):
    """
    Log logs a message to the console.
    """
    sys.stderr.write(msg + '\n')

if __name__ == '__main__':
    main()
