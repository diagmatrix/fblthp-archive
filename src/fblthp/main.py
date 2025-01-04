from scryfall import Scryfall
import sqlite3

if __name__ == "__main__":
    connection = sqlite3.connect("../../archive.db")
    client = Scryfall(connection)

    sets = client.get_sets()
    print(sets)