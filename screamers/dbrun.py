from config import TG_TOKEN
from config import MONGODB_LINK
from config import MONGO_DB
from pymongo import MongoClient
import time

mondb = MongoClient(MONGODB_LINK)[MONGO_DB]

resume_token = None
pipeline = [{'$match': {'operationType': 'update'}}]
with mondb.runners.watch(pipeline=pipeline) as stream:
    while stream.alive:
        change = stream.try_next()
        # Note that the ChangeStream's resume token may be updated
        # even when no changes are returned.
        if change is not None:
            print("Change document: %r" % (change,))
            continue
        # We end up here when there are no recent changes.
        # Sleep for a while before trying again to avoid flooding
        # the server with getMore requests when no changes are
        # available.
        time.sleep(10)