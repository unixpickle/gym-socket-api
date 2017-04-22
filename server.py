import os
import socket
import traceback
import proto
import gym
import numpy as np
import json

def serve(port):
    # TODO: use socketserver module.
    s = socket.socket()
    s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    s.bind(('127.0.0.1', port))
    print('Listening on port %d...' % port)
    s.listen(10)
    while True:
        conn, addr = s.accept()
        pid = os.fork()
        if pid == 0:
            handle(conn, addr)
            exit()
        else:
            conn.close()

def handle(conn, addr):
    print('Connection from %s' % str(addr))
    f = conn.makefile(mode='rwb')
    try:
        env = handshake(f)
        try:
            loop(f, env)
        finally:
            env.close()
    except Exception as e:
        print('Error from %s' % str(addr))
        traceback.print_exc()
    finally:
        f.close()

def handshake(sock):
    flags = proto.read_flags(sock)
    envName = proto.read_field_str(sock)
    try:
        env = gym.make(envName)
        proto.write_field_str(sock, '')
        sock.flush()
        return env
    except gym.error.Error as e:
        proto.write_field_str(sock, str(e))
        sock.flush()
        raise e

def loop(sock, env):
    while True:
        pack_type = proto.read_packet_type(sock)
        if pack_type == 'reset':
            handle_reset(sock, env)
        elif pack_type == 'step':
            handle_step(sock, env)

def handle_reset(sock, env):
    send_obs(sock, env, env.reset())
    sock.flush()

def handle_step(sock, env):
    action = proto.read_action(sock)
    if isinstance(action, list):
        action = np.array(action)
    obs, rew, done, info = env.step(action)
    send_obs(sock, env, obs)
    proto.write_reward(sock, rew)
    proto.write_bool(sock, done)
    proto.write_field_str(sock, json.dumps(info))
    sock.flush()

def send_obs(sock, env, obs):
    if isinstance(obs, np.ndarray):
        if obs.dtype == 'uint8':
            proto.write_obs_byte_list(sock, obs)
            return
    jsonable = env.observation_space.to_jsonable(obs)
    proto.write_obs_json(sock, jsonable)
