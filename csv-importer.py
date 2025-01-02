import os
import sqlite3
import pandas as pd

def main():
    # Config
    TABLE_NAME = "raw_collection"
    BASE_QUERY = f"INSERT INTO {TABLE_NAME} "

    # Connect to SQLite
    conn = sqlite3.connect('archive.db')
    c = conn.cursor()

    # Get all csv files from data folder
    import_files = []
    for root, dirs, files in os.walk("data"):
        for file in files:
            if file.endswith(".csv"):
                import_files.append(os.path.join(root, file))

    for file in import_files:
        # Import files
        df = pd.read_csv(file, sep=";")
        print(f"Added {file} to memory ({len(df)} rows)")

        # Create query
        cols = ", ".join([column_parse(col) for col in df.columns])
        values = ", ".join("?"*len(df.columns))
        query = f"{BASE_QUERY} ({cols}) VALUES ({values})"
        print(f"Executing query: \"{query}\"")

        # Add to database
        cards = df.values.tolist()
        c.executemany(query, cards)
        conn.commit()
        print(f"Cards successfully imported")

    conn.commit()
    conn.close()

def column_parse(col_name: str):
    col = col_name.lower()
    # Check if set
    if col == "set":
        col = f"\"{col}\""
    # Check if space
    if " " in col:
        col = col.replace(" ", "_")

    return col

if __name__ == "__main__":
    main()