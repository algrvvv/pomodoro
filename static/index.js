document.addEventListener('DOMContentLoaded', () => {
  fetch("/api/v1/data")
    .then(response => response.json())
    .then(data => {
      console.log(data)
      renderCalendar(data.calendar);
      // renderChart(data.chart)
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

function renderCalendar(dates) {
  const calendar = document.querySelector(".calendar");
  calendar.innerHTML = '<h3>Календарь сессий</h3>';

  const daysContainer = document.createElement("div");
  daysContainer.className = "calendar-days";
  calendar.appendChild(daysContainer);

  // Получаем текущий месяц и год
  const now = new Date();
  const year = now.getFullYear();
  const month = now.getMonth(); // 0-индексация для месяцев: 0 = январь

  // Следующий месяц с днем 0 = последний день текущего месяца
  const daysInMonth = new Date(year, month + 1, 0).getDate(); 

  for (let i = 1; i <= daysInMonth; i++) {
    const day = document.createElement("span");
    day.className = "day";
    day.textContent = i;

    // Формируем дату в формате YYYY-MM-DD
    const date = `${year}-${String(month + 1).padStart(2, "0")}-${String(i).padStart(2, "0")}`;

    // Проверяем, есть ли дата в полученных данных
    if (dates.includes(date)) {
      day.classList.add("selected"); // Подсвечиваем дни, которые есть в данных
    }

    daysContainer.appendChild(day); // Добавляем день в контейнер
  }
}

function renderChart(chartData) {
  const chartContainer = document.querySelector(".chart-bars");
  chartContainer.innerHTML = ""; // Очищаем график перед рендерингом

  // Найдем максимальное значение для масштабирования
  const maxSessions = Math.max(...Object.values(chartData));

  // Преобразуем данные в массив [дата, значение] и отсортируем по дате
  const sortedData = Object.entries(chartData).sort(([dateA], [dateB]) => new Date(dateA) - new Date(dateB));

  // Создаем бары для графика
  sortedData.forEach(([date, count]) => {
    const bar = document.createElement("div");
    bar.className = "bar";
    bar.style.height = `${(count / maxSessions) * 100}%`; // Высота в процентах от максимального значения
    bar.title = `${date}: ${count} сессий`; // Подсказка при наведении
    chartContainer.appendChild(bar);

    // Добавляем подпись под баром
    const label = document.createElement("span");
    label.className = "label";
    label.textContent = new Date(date).getDate(); // День месяца
    chartContainer.appendChild(label);
  });
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
          borderColor: "#4caf50", // Цвет линии
          backgroundColor: "rgba(76, 175, 80, 0.2)", // Полупрозрачный цвет заливки
          borderWidth: 2,
          pointBackgroundColor: "#4caf50", // Цвет точек
          pointRadius: 4, // Размер точек
          fill: true, // Включить заливку под линией
          tension: 0.4, // Сделать линию сглаженной
        },
      ],
    },
    options: {
      responsive: true,
      maintainAspectRatio: true, // включить, если размер родителя не зафиксирован 
      interaction: {
        mode: "nearest", // Отображать подсказку для ближайшей точки
        axis: "x", // Ограничиваем по оси X
        intersect: false, // Не обязательно пересекать точку мышью
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

