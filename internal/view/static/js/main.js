// Функция для переключения видимости выпадающего меню
function toggleDropdown() {
    const dropdown = document.getElementById('dropdown-menu');
    dropdown.style.display = dropdown.style.display === 'flex' ? 'none' : 'flex';
}

// Закрытие выпадающего меню при клике вне его
window.onclick = function(event) {
    const profileContainer = document.querySelector('.profile-container');
    const dropdown = document.getElementById('dropdown-menu');
    if (!profileContainer.contains(event.target) && !dropdown.contains(event.target)) {
        dropdown.style.display = 'none';
    }
};
