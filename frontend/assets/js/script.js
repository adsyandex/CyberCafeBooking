// Обработка формы бронирования
document.getElementById('bookingForm').addEventListener('submit', async (e) => {
    e.preventDefault();

    // Сбор данных из формы
    const booking = {
        name: document.getElementById('name').value,
        contact: document.getElementById('contact').value,
        computer: document.getElementById('computer').value,
        startTime: document.getElementById('startTime').value,
        endTime: document.getElementById('endTime').value
    };

    // Отправка данных на сервер
    const response = await fetch('/api/bookings', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(booking)
    });

    const data = await response.json(); // Получаем обновленный список
    updateTable(data);
});

// Загрузка существующих бронирований
async function loadBookings() {
    const response = await fetch('/api/bookings');
    const data = await response.json();
    updateTable(data);
}

// Обновление таблицы бронирований
function updateTable(bookings) {
    const tableBody = document.getElementById('bookingTable');
    tableBody.innerHTML = ''; // Очистка таблицы
    bookings.forEach(booking => {
        const row = `<tr>
            <td>${booking.name}</td>
            <td>${booking.contact}</td>
            <td>${booking.computer}</td>
            <td>${booking.startTime}</td>
            <td>${booking.endTime}</td>
        </tr>`;
        tableBody.innerHTML += row;
    });
}

loadBookings(); // Загрузка данных при открытии страницы
