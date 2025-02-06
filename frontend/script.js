const ctx = document.getElementById('performanceChart').getContext('2d');
let chart;

function fetchData() {
  fetch('http://localhost:8080/performance')
    .then(response => response.json())
    .then(data => {
      if (!chart) {
        // Инициализация графика
        chart = new Chart(ctx, {
          type: 'line',
          data: {
            labels: data.map(entry => entry.workers),
            datasets: [{
              label: 'Time (seconds)',
              data: data.map(entry => entry.time),
              borderColor: 'rgba(75, 192, 192, 1)',
              borderWidth: 2,
              fill: false,
            }]
          },
          options: {
            scales: {
              x: {
                title: {
                  display: true,
                  text: 'Number of Workers'
                }
              },
              y: {
                title: {
                  display: true,
                  text: 'Time (seconds)'
                }
              }
            }
          }
        });
      } else {
        // Обновление данных графика
        chart.data.labels = data.map(entry => entry.workers);
        chart.data.datasets[0].data = data.map(entry => entry.time);
        chart.update();
      }
    })
    .catch(error => console.error('Error fetching data:', error));
}

// Обновление данных каждые 5 секунд
setInterval(fetchData, 5000);

// Первый запрос данных
fetchData();