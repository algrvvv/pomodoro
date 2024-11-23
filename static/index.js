document.addEventListener('DOMContentLoaded', () => {
  const stopwatch = document.getElementById('stats-stopwatch')
  const socket = new WebSocket(`ws://localhost:3333/ws`)

  socket.onopen = () => {
    console.log('websocket connected')
    setTimeout(() => {
      // console.log(stopwatch.innerHTML)
      if (stopwatch.innerHTML == "Загрузка...") {
        stopwatch.innerHTML = "Нет активной сессии"
      }
    }, 3000)
  }

  socket.onmessage = (msg) => {
    stopwatch.innerHTML = msg.data
    // console.log('got data: ', msg.data)
  }

  socket.onclose = () => {
    console.warn('websocket closed')
  }

  socket.onerror = (err) => {
    console.error(err)
  }

  fetch("/api/v1/data")
    .then(response => response.json())
    .then(data => {
      console.log(data)

      const today = new Date();
      const todayString = `${today.getFullYear()}-${today.getMonth() + 1}-${today.getDate()}`

      renderCalendar(today.getMonth(), today.getFullYear(), data.calendar);
      renderSessions(data.sessions);

      renderLineChart(data.chart, todayString)

      renderTotalMinutes(data.totalMinutes)
    });
});

function renderTotalMinutes(minutes) {
  document.getElementById('stats-total-minutes').innerHTML = minutes
}

function renderSessions(sessions) {
  // console.log('got sessions: ', sessions)
  const container = document.querySelector('.sessions ul')
  container.innerHTML = ""

  if (!sessions) { return }

  let totalToday = 0
  sessions.forEach(s => {
    totalToday += s.duration

    const li = document.createElement('li')
    const type = s.type == 1 ? 'Работа' : 'Отдых'

    const getCorrectFormat = function(time) {
      const minutes = time.getMinutes().toString()
      const hours = time.getHours().toString()
      return `${hours.padStart(2, '0')}:${minutes.padStart(2, '0')}`
    }

    const endDate = new Date(s.date)
    const startDate = new Date(s.date).setMinutes(endDate.getMinutes() - s.duration)
    const startTime = new Date(startDate)

    const fStartTime = getCorrectFormat(startTime)
    const fEndTime = getCorrectFormat(endDate)

    const time = `${fStartTime} - ${fEndTime}`
    const date = `${endDate.getDate()}.${endDate.getMonth()}.${endDate.getFullYear()} ${time}`

    li.textContent = `Сессия ${s.id} (${date}): Продолжительность ${s.duration} мин - ${type}`
    container.appendChild(li)
  });

  document.getElementById('stats-today-minutes').innerHTML = totalToday;
}

function renderCalendar(month, year, dates) {
  const calendar = document.querySelector(".calendar");
  calendar.innerHTML = '<h3>Достижение цели</h3>';

  const daysContainer = document.createElement("div");
  daysContainer.className = "calendar-days";
  calendar.appendChild(daysContainer);

  const daysInMonth = new Date(year, month + 1, 0).getDate();
  const firstDayOfWeek = new Date(year, month, 1).getDay();
  const correctedFirstDay = firstDayOfWeek === 0 ? 6 : firstDayOfWeek - 1;

  for (let i = 0; i < correctedFirstDay; i++) {
    const emptyCell = document.createElement("div");
    emptyCell.classList.add("empty");
    daysContainer.appendChild(emptyCell);
  }

  let streak = 0;
  const today = new Date().getDate();
  for (let i = 1; i <= daysInMonth; i++) {
    const day = document.createElement("span");
    day.className = "day";
    day.textContent = i;

    // форматируем дату в формате YYYY-MM-DD
    const date = `${year}-${String(month + 1).padStart(2, "0")}-${String(i).padStart(2, "0")}`;

    if (dates.includes(date)) {
      day.classList.add("selected");
      streak++
    } else {
      if (today >= i) { streak = 0 }
    }

    daysContainer.appendChild(day);
  }
  document.getElementById('streak').innerHTML = streak;
}

function renderLineChart(chartData, today) {
  if (!chartData[today]) { chartData[today] = 0 }

  const ctx = document.getElementById("lineChart").getContext("2d");

  const sortedData = Object.entries(chartData).sort(([dateA], [dateB]) => new Date(dateA) - new Date(dateB));

  const labels = sortedData.map(([date]) => new Date(date).toLocaleDateString());
  const data = sortedData.map(([, value]) => value);

  new Chart(ctx, {
    type: "line",
    data: {
      labels: labels,
      datasets: [
        {
          label: "Количество сессий",
          data: data,
          borderColor: "#4caf50",
          backgroundColor: "rgba(76, 175, 80, 0.2)",
          borderWidth: 2,
          pointBackgroundColor: "#4caf50",
          pointRadius: 2,
          fill: true,
          tension: 0.4,
        },
      ],
    },
    options: {
      responsive: true,
      maintainAspectRatio: true,
      interaction: {
        mode: "nearest",
        axis: "x",
        intersect: false,
      },
      plugins: {
        tooltip: {
          enabled: true,
          callbacks: {
            label: function(context) {
              return `Дата: ${context.label}, Сессии: ${context.raw}`;
            },
          },
        },
        legend: {
          display: true,
        },
      },
      scales: {
        x: {
          title: {
            display: true,
            text: "Дата",
          },
        },
        y: {
          title: {
            display: true,
            text: "Количество сессий",
          },
          beginAtZero: true,
        },
      },
    },
  });
}

