from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from datetime import datetime
import psycopg2
import os

app = FastAPI(title="LIM Analytics")


def get_conn():
    return psycopg2.connect(
        host=os.getenv("DB_HOST", "postgres-master"),
        dbname=os.getenv("DB_NAME", "postgres"),
        user=os.getenv("DB_USER", "postgres"),
        password=os.getenv("DB_PASSWORD", "12345"),
        port=int(os.getenv("DB_PORT", "5432")),
    )


class Stats(BaseModel):
    clients: int
    coaches: int
    halls: int
    trainings: int


class CoachStats(BaseModel):
    name: str
    trainings: int
    clients: int
    total_places: int
    fill_rate: float
    styles: str


class Training(BaseModel):
    id: int
    name: str
    date_time: datetime
    places: int
    coach: str
    hall: int
    signed_up: int


@app.get("/stats", response_model=Stats)
def get_stats():
    conn = get_conn()
    cur = conn.cursor()
    cur.execute("SELECT COUNT(*) FROM clients")
    clients = cur.fetchone()[0]
    cur.execute("SELECT COUNT(*) FROM coaches")
    coaches = cur.fetchone()[0]
    cur.execute("SELECT COUNT(*) FROM halls")
    halls = cur.fetchone()[0]
    cur.execute("SELECT COUNT(*) FROM trainings")
    trainings = cur.fetchone()[0]
    cur.close()
    conn.close()
    return Stats(clients=clients, coaches=coaches, halls=halls, trainings=trainings)


@app.get("/coaches/stats", response_model=list[CoachStats])
def get_coaches_stats():
    conn = get_conn()
    cur = conn.cursor()
    cur.execute("""
        SELECT
            c.name,
            COUNT(DISTINCT t.training_id) AS trainings,
            COUNT(ct.client_id) AS clients,
            COALESCE(SUM(t.places_num), 0) AS total_places,
            COALESCE(
                ROUND(COUNT(ct.client_id)::numeric / NULLIF(SUM(t.places_num), 0) * 100, 1),
                0
            ) AS fill_rate,
            STRING_AGG(DISTINCT t.name, ', ' ORDER BY t.name) AS styles
        FROM coaches c
        LEFT JOIN trainings t ON c.coach_id = t.coach_id
        LEFT JOIN clients_trainings ct ON t.training_id = ct.training_id
        GROUP BY c.coach_id, c.name
        ORDER BY trainings DESC
    """)
    rows = cur.fetchall()
    cur.close()
    conn.close()
    return [
        CoachStats(
            name=r[0],
            trainings=r[1],
            clients=r[2],
            total_places=r[3],
            fill_rate=float(r[4]),
            styles=r[5] or "—",
        )
        for r in rows
    ]


@app.get("/trainings", response_model=list[Training])
def get_trainings():
    conn = get_conn()
    cur = conn.cursor()
    cur.execute("""
        SELECT
            t.training_id,
            t.name,
            t.date_time,
            t.places_num,
            c.name,
            h.number,
            COUNT(ct.client_id) AS signed_up
        FROM trainings t
        JOIN coaches c ON t.coach_id = c.coach_id
        JOIN halls h ON t.hall_id = h.hall_id
        LEFT JOIN clients_trainings ct ON t.training_id = ct.training_id
        GROUP BY t.training_id, c.name, h.number
        ORDER BY t.date_time
    """)
    rows = cur.fetchall()
    cur.close()
    conn.close()
    return [
        Training(id=r[0], name=r[1], date_time=r[2], places=r[3], coach=r[4], hall=r[5], signed_up=r[6])
        for r in rows
    ]
