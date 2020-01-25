"""
extractor.py : This program is meant to observe the
incoming beat files and extract the metadata from it.
"""

import time
import logging
import os
import threading

from watchdog.observers import Observer
from watchdog.events import LoggingEventHandler
from watchdog.events import FileSystemEventHandler

from lib.metadata import get_beats, get_verse_intervals, get_bpm

POLL_DELAY = os.environ.get('POLL_DELAY')
WATCH_PATH = os.environ.get('WATCH_PATH')
BEAT_EXTENSION = os.environ.get('BEAT_EXTENSION') # mp3

class FSHandler(FileSystemEventHandler):
    """
    The operation ran when a new file is created into the
    watched filesystem.
    """
    def on_created(self, event):
        filename = event.src_path.split("/")[-1]
        if filename.split(".")[1] == BEAT_EXTENSION:
            print("[EXTRACTOR] Starting metadata extraction of : "+filename)
            threading.Thread(target=self.__extract, args=(filename,)).start()
        else:
            print("[EXTRACTOR] The file : "+filename+" has the wrong extension. (It should be ."+BEAT_EXTENSION+")")

    def __extract(self, filename):
        get_bpm(WATCH_PATH+filename)
        beats_dist, beatfile_duration = get_beats(WATCH_PATH+filename)
        get_verse_intervals(WATCH_PATH+filename, beats_dist, beatfile_duration)
        print("[EXTRACTOR] Metadata extracted successfully !")

if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO,
                        format='%(asctime)s - %(message)s',
                        datefmt='%Y-%m-%d %H:%M:%S')
    LOGGING_EVENT_HANDLER = LoggingEventHandler()
    FILESYS_EVENT_HANDLER = FSHandler()
    OBSERVER = Observer()
    OBSERVER.schedule(LOGGING_EVENT_HANDLER, WATCH_PATH, recursive=False)
    OBSERVER.schedule(FILESYS_EVENT_HANDLER, WATCH_PATH, recursive=False)
    OBSERVER.start()
    try:
        while True:
            time.sleep(POLL_DELAY)
    except KeyboardInterrupt:
        OBSERVER.stop()
    OBSERVER.join()
