document.addEventListener('DOMContentLoaded', () => {
  fetch("/api/v1/data")
    .then(response => response.json())
    .then(data => {
      console.log(data)

      const today = new Date();
      renderCalendar(today.getMonth(), today.getFullYear(), data.calendar);
      renderSessions(data.sessions);
      renderLineChart(data.chart)
    });
});

function renderSessions(sessions) {
  console.log('got sessions: ', sessions)
  const container = document.querySelector('.sessions ul')
  container.innerHTML = ""

  sessions.forEach(s => {
    const li = document.createElement('li')
    const type = s.type == 1 ? 'Работа' : 'Отдых'
    li.textContent = `Сессия ${s.id}: Продолжительность ${s.duration} мин - ${type}`
    container.appendChild(li)
  });
}

function renderCalendar(month, year, dates) {
  const calendar = document.querySelector(".calendar");
  calendar.innerHTML = '<h3>Календарь сессий</h3>';

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

  for (let i = 1; i <= daysInMonth; i++) {
    const day = document.createElement("span");
    day.className = "day";
    day.textContent = i;

    // форматируем дату в формате YYYY-MM-DD
    const date = `${year}-${String(month + 1).padStart(2, "0")}-${String(i).padStart(2, "0")}`;

    if (dates.includes(date)) {
      day.classList.add("selected");
    }

    daysContainer.appendChild(day);
  }
}

function renderLineChart(chartData) {
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
          pointRadius: 4,
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
            label: function (context) {
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

