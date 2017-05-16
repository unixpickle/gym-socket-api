"""
APIs for Universe-specific environment manipulations.
"""

import numpy as np
import gym

class UniverseException(Exception):
    """
    Exception type used for all Universe-related errors.
    """
    pass

class Universe:
    """
    Universe wraps access to OpenAI Universe.
    """
    def __init__(self, enabled):
        self.enabled = enabled
        if enabled:
            import universe
            self.universe = universe

    def wrap(self, env, wrapper_name, options):
        """
        Wrap wraps the environment.

        This should be called before the environment is
        configured.
        """
        self._check_enabled()
        self._check_env(env)
        wrappers = self.universe.wrappers
        classes = {
            'CropObservations': wrappers.experimental.CropObservations,
            'Vision': wrappers.Vision
        }
        if not wrapper_name in classes:
            raise UniverseException('unknown wrapper: ' + wrapper_name)
        try:
            return classes[wrapper_name](env, **options)
        except TypeError as exc:
            raise UniverseException('bad wrapper options: ' + str(exc))

    def configure(self, env, options):
        """
        Configure a Universe environment.

        Returns a new (potentially wrapped) environment.
        """
        self._check_enabled()
        self._check_env(env)
        try:
            env.configure(**options)
        except TypeError as exc:
            raise UniverseException('bad configure options: ' + str(exc))
        wrappers = self.universe.wrappers
        return UniverseEnv(wrappers.Unvectorize(wrappers.BlockingReset(env)))

    def _check_enabled(self):
        if not self.enabled:
            raise UniverseException('Universe is not enabled')

    def _check_env(self, env):
        if not isinstance(env, gym.Env):
            raise UniverseException('not a Gym environment')
        if not isinstance(env.unwrapped, self.universe.envs.VNCEnv):
            raise UniverseException('not a Universe environment')

class UniverseEnv:
    """
    A pseudo environment wrapper with more useful spaces.
    """
    def __init__(self, env):
        self.env = env
        self.action_space = UniverseActionSpace(env.action_space)
        self.observation_space = UniverseObservationSpace(env.observation_space)

    def reset(self):
        """
        Reset the environment.
        """
        return self.env.reset()

    def step(self, action):
        """
        Take a step in the environment.
        """
        return self.env.step(action)

class UniverseActionSpace:
    """
    A pseudo action space for VNC environments.

    This is necessary because the real action space does
    not support JSON.
    """
    def __init__(self, space):
        self.space = space

    @staticmethod
    def sample():
        """
        Sample from the space.
        """
        raise UniverseException('unable to sample VNC actions')

    @staticmethod
    def to_jsonable(sample_n):
        """
        Convert a batch of actions to JSON.
        """
        return [[list(t) for t in sample] for sample in sample_n]

    @staticmethod
    def from_jsonable(sample_n):
        """
        Convert a batch of actions from JSON.
        """
        return [[tuple(t) for t in sample] for sample in sample_n]

# pylint: disable=too-few-public-methods
class UniverseObservationSpace:
    """
    A pseudo observation space for VNC environments.

    Just like UniverseActionSpace, this is necessary
    because of JSON.
    """
    def __init__(self, space):
        self.space = space

    @classmethod
    def to_jsonable(cls, sample_n):
        """
        Convert a batch of observations to JSON.
        """
        return [cls._to_jsonable(obj) for obj in sample_n]

    @classmethod
    def _to_jsonable(cls, obj):
        if isinstance(obj, np.ndarray):
            return obj.tolist()
        elif isinstance(obj, dict):
            res = {}
            for key in obj:
                res[key] = cls._to_jsonable(obj[key])
            return res
        return obj
