class MockOptions(object):
    """Simple object for storing attributes.
    """

    def __init__(self, **kwargs):
        for name in kwargs:
            setattr(self, name, kwargs[name])
