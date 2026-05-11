import streamlit as st
import requests
import pandas as pd
from datetime import datetime

API_URL = "http://lim-analytics:8000"

st.title("Lost in Motion — аналитика студии")

stats_resp = requests.get(f"{API_URL}/stats")
if stats_resp.status_code == 200:
    stats = stats_resp.json()
    col1, col2, col3, col4 = st.columns(4)
    col1.metric("Клиенты", stats["clients"])
    col2.metric("Тренеры", stats["coaches"])
    col3.metric("Залы", stats["halls"])
    col4.metric("Тренировки", stats["trainings"])

st.divider()
st.subheader("Статистика по тренерам")

coaches_resp = requests.get(f"{API_URL}/coaches/stats")
if coaches_resp.status_code != 200:
    st.error("Не удалось загрузить данные")
    st.stop()

coaches = coaches_resp.json()
df_coaches = pd.DataFrame(coaches).rename(columns={
    "name": "Тренер",
    "trainings": "Тренировок",
    "clients": "Записей",
    "total_places": "Всего мест",
    "fill_rate": "Заполненность, %",
    "styles": "Направления",
})
st.dataframe(df_coaches, use_container_width=True, hide_index=True)

st.divider()
st.subheader("Расписание тренера")

selected = st.selectbox("Выбери тренера", [c["name"] for c in coaches])

trainings_resp = requests.get(f"{API_URL}/trainings")
if trainings_resp.status_code != 200:
    st.error("Не удалось загрузить тренировки")
    st.stop()

all_trainings = trainings_resp.json()
coach_trainings = [t for t in all_trainings if t["coach"] == selected]

if not coach_trainings:
    st.info("У этого тренера нет тренировок")
    st.stop()

tdf = pd.DataFrame(coach_trainings)
tdf["date_time"] = pd.to_datetime(tdf["date_time"])
tdf["свободно"] = tdf["places"] - tdf["signed_up"]

now = datetime.now()
upcoming = tdf[tdf["date_time"] > now].sort_values("date_time")
past = tdf[tdf["date_time"] <= now].sort_values("date_time", ascending=False)

tab1, tab2 = st.tabs(["Предстоящие", "Проведённые"])

def format_df(df):
    df = df.copy()
    df["date_time"] = df["date_time"].dt.strftime("%d.%m.%Y %H:%M")
    return df[["name", "date_time", "hall", "places", "signed_up", "свободно"]].rename(columns={
        "name": "Направление",
        "date_time": "Дата и время",
        "hall": "Зал №",
        "places": "Мест всего",
        "signed_up": "Записалось",
    })

with tab1:
    if upcoming.empty:
        st.info("Нет предстоящих тренировок")
    else:
        st.dataframe(format_df(upcoming), use_container_width=True, hide_index=True)

with tab2:
    if past.empty:
        st.info("Нет проведённых тренировок")
    else:
        st.dataframe(format_df(past), use_container_width=True, hide_index=True)
