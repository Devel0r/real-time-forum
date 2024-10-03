document.addEventListener('DOMContentLoaded', function() {
    const profile = document.querySelector('.profile');
    const dropdownMenu = document.querySelector('.dropdown-menu');

    profile.addEventListener('click', function(e) {
        e.stopPropagation();
        dropdownMenu.classList.toggle('show');
    });

    // Закрытие выпадающего меню при клике вне его
    document.addEventListener('click', function(e) {
        if (!profile.contains(e.target)) {
            dropdownMenu.classList.remove('show');
        }
    });
});
