document.addEventListener('DOMContentLoaded', () => {
  const stopwatch = document.getElementById('stats-stopwatch')
  const socket = new WebSocket(`ws://localhost:3333/ws`)
  const monthCalendar = document.getElementById('calendar')

  socket.onopen = () => {
    console.log('websocket connected')
    setTimeout(() => {
      if (stopwatch.innerHTML.includes('svg')) {
        stopwatch.innerHTML = "<span>Нет активной сессии</span>"
      }
    }, 3000)
  }

  socket.onmessage = (msg) => {
    stopwatch.innerHTML = msg.data
  }

  socket.onclose = () => {
    console.warn('websocket closed')
  }

  socket.onerror = (err) => {
    console.error(err)
  }


  // date format: "2024-12-01"
  let date = localStorage.getItem("date")
  if (!date) {
    const d = new Date();
    date = `${d.getFullYear()}-${d.getMonth() + 1}-01`
    localStorage.setItem("date", date)
  }
  let splitedDate = date.split('-')
  const dateWithoutDay = `${splitedDate[0]}-${splitedDate[1]}`

  monthCalendar.value = dateWithoutDay
  monthCalendar.addEventListener('change', () => {
    let newDate = `${monthCalendar.value}-01`
    localStorage.setItem("date", newDate)
    window.location.reload()
  })

  fetch(`/api/v1/data?date=${date}`)
    .then(response => response.json())
    .then(data => {
      console.log("got data from request: ", data)
      document.getElementById('time').textContent = `Время загрузки: ${data.time}`

      data.integrations.forEach(i => {
        if (i.name.toLowerCase() == "wakatime" && i.enabled) {

          changeVisible('wakatime-integration', 'block')
          console.log('start fetching data...')
          fetchWakatimeStats().then(data => {
            changeVisible('wakatime-loading', 'none')

            console.log("данные по статистике: ", data)
            document.getElementById('wakatime-today').textContent = data.today.data
            document.getElementById('wakatime-week').textContent = data.week.data
            document.getElementById('wakatime-loadtime').textContent = data.time

            changeVisible('wakatime-content', 'flex')
          }).catch(err => console.error("err: ", err))
        }
      })
      console.log('continue render...')

      const today = new Date(dateWithoutDay);
      // const todayDate = today.getDate().toString().padStart(2, '0')
      // const todayString = `${today.getFullYear()}-${today.getMonth() + 1}-${todayDate}`

      const todayForChart = new Date();
      let todayForChartString;
      if (todayForChart.getMonth() == today.getMonth()) {
        todayForChartString = formatDate(todayForChart)
      } else {
        todayForChartString = formatDate(today)
      }

      renderCalendar(today.getMonth(), today.getFullYear(), data.calendar, data.tooltips);
      renderSessions(data.sessions);

      renderLineChart(data.chartCount, data.chartMinutes, todayForChartString)

      renderTotalMinutes(data.totalMinutes)
      addMouseOverEvent()
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
    const date = `${endDate.getDate()}.${endDate.getMonth() + 1}.${endDate.getFullYear()} ${time}`

    li.textContent = `Сессия ${s.id}(${date}): Продолжительность ${s.duration} мин - ${type}`
    container.appendChild(li)
  });

  document.getElementById('stats-today-minutes').innerHTML = totalToday;
}

function renderCalendar(month, year, dates, tooltips) {
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
    let dataForTooltip = "Сессий не было"
    if (tooltips[i]) {
      const t = tooltips[i]
      dataForTooltip = `Статистика за ${i} число: <br><br>
      ${t[1] ? 'Работа: ' + t[1] + ' мин.<br>' : ''}
      ${t[2] ? 'Отдых: ' + t[2] + ' мин.' : ''}
      `
    }

    const day = document.createElement("span");
    day.className = "day";
    day.setAttribute("data", dataForTooltip)
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

function renderLineChart(chartData, chartDataMinutes, today) {
  if (!chartData[today]) { chartData[today] = 0 }
  if (!chartDataMinutes[today]) { chartDataMinutes[today] = 0 }

  const ctx = document.getElementById("lineChart").getContext("2d");

  const sortedData = Object.entries(chartData).sort(([dateA], [dateB]) => new Date(dateA) - new Date(dateB));
  const sortedMinutesData = Object.entries(chartDataMinutes).sort(([dateA], [dateB]) => new Date(dateA) - new Date(dateB));

  const labels = sortedData.map(([date]) => new Date(date).toLocaleDateString());
  const data = sortedData.map(([, value]) => value);
  const minutesData = sortedMinutesData.map(([, value]) => value);

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
          yAxisID: 'y1', // Используем первую ось Y
        },
        {
          label: "Количество минут",
          data: minutesData,
          borderColor: "#8b5cf6",
          backgroundColor: "rgba(196, 181, 253, 0.2)",
          borderWidth: 2,
          pointBackgroundColor: "#c4b5fd",
          pointRadius: 2,
          fill: true,
          tension: 0.4,
          yAxisID: 'y2', // Используем первую ось Y
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
              const datasetLabel = context.dataset.label;
              const value = context.raw;
              const label = context.label;

              if (datasetLabel === "Количество сессий") {
                return `Дата: ${label}, Сессии: ${value}`;
              }

              // Если набор данных - "Количество минут"
              if (datasetLabel === "Количество минут") {
                return `Дата: ${label}, Минут: ${value}`;
              }

              // Для других случаев (если нужно обработать)
              return `${datasetLabel}: ${value}`;
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
        y1: {
          title: {
            display: true,
            text: "Количество сессий",
          },
          // ticks: { stepSize: 5 },
          grid: { display: true },
          beginAtZero: true,
        },
        y2: {
          title: {
            display: true,
            text: "Количество минут",
          },
          // выключаем отображение сетки для второго игрика,
          // чтобы лучше выглядело
          grid: { display: false },
          beginAtZero: true,
          position: 'right'
        }
      },
    },
  });
}

function addMouseOverEvent() {
  const blocks = document.querySelectorAll('.day')
  const tooltip = document.getElementById('tooltip')

  blocks.forEach(b => {
    b.addEventListener('mouseover', (event) => {
      const data = b.getAttribute("data")

      tooltip.innerHTML = data;
      tooltip.style.top = event.pageY + 10 + 'px';
      tooltip.style.left = event.pageX + 10 + 'px';
      tooltip.classList.add('visible');
    })

    b.addEventListener('mousemove', (event) => {
      tooltip.style.top = event.pageY + 10 + 'px';
      tooltip.style.left = event.pageX + 10 + 'px';
    })

    b.addEventListener('mouseout', () => {
      tooltip.classList.remove('visible')
    })
  })
}

async function fetchWakatimeStats() {
  const url = `/api/v1/integrations/wakatime`
  const response = await fetch(url, {
    headers: {
      'Content-Type': 'application/json'
    }
  });

  if (!response.ok) {
    throw new Error(`Ошибка: ${response.status}`);
  }
  const data = await response.json();
  return data;
}

function formatDate(date) {
  return `${date.getFullYear()}-${date.getMonth() + 1}-${date.getDate().toString().padStart(2, '0')}`
}

function changeVisible(id, mode) {
  document.getElementById(id).style.display = mode;
}

