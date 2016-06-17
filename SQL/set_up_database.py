import psycopg2
import os

def populate_data():
    # Clean up and recreate the database. Database is populated with some sample data
    conn = psycopg2.connect(
        database=os.environ['PROJ_DB_NAME'],
        user=os.environ['PROJ_DB_USER'],
        host=os.environ['PROJ_DB_HOST'],
        password=os.environ['PROJ_DB_PWD'],
        port=int(os.environ['PROJ_DB_PORT'])
    )
    cursor = conn.cursor()
    cursor.execute(open("01_setting_up.sql", "r").read())
    cursor.execute(open("02_populate.sql", "r").read())
    conn.commit()

populate_data()