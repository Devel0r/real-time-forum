// Получаем ссылки на элементы
const loginForm = document.getElementById('login-form');
const signupForm = document.getElementById('signup-form');
const postsSection = document.getElementById('posts-section');
const chatSection = document.getElementById('chat-section');
const authSection = document.getElementById('auth-section');
const appSection = document.getElementById('app-section');
const loginBtn = document.getElementById('login-btn');
const signupBtn = document.getElementById('signup-btn');
const postsBtn = document.getElementById('posts-btn');
const chatBtn = document.getElementById('chat-btn');
const logoutBtn = document.getElementById('logout-btn');

// Функции для показа/скрытия секций
function showLoginForm() {
    authSection.style.display = 'block';
    loginForm.style.display = 'block';
    signupForm.style.display = 'none';
    appSection.style.display = 'none';
}

function showSignupForm() {
    authSection.style.display = 'block';
    signupForm.style.display = 'block';
    loginForm.style.display = 'none';
    appSection.style.display = 'none';
}

function showAppSection() {
    authSection.style.display = 'none';
    appSection.style.display = 'block';
    postsSection.style.display = 'block';
    chatSection.style.display = 'none';
}

function showPostsSection() {
    postsSection.style.display = 'block';
    chatSection.style.display = 'none';
}

function showChatSection() {
    chatSection.style.display = 'block';
    postsSection.style.display = 'none';
}

// Слушатели событий для кнопок навигации
postsBtn.addEventListener('click', showPostsSection);
chatBtn.addEventListener('click', showChatSection);
logoutBtn.addEventListener('click', () => {
    // Логика выхода (удаление сессии на сервере)
    fetch('/api/logout', {
        method: 'POST'
    })
        .then(response => {
            if (response.ok) {
                alert('Logged out!');
                showLoginForm();
            } else {
                alert('Error logging out');
            }
        })
        .catch(error => {
            console.error('Error:', error);
        });
});

// Обработчик отправки формы входа
document.getElementById('login').addEventListener('submit', function(e) {
    e.preventDefault();

    const formData = new FormData(e.target);

    fetch('/api/login', {
        method: 'POST',
        body: formData
    })
        .then(response => {
            if (response.ok) {
                return response.json();
            } else {
                return response.json().then(errorData => {
                    throw new Error(errorData.message || 'Login failed');
                });
            }
        })
        .then(data => {
            alert(data.message);
            showAppSection();
        })
        .catch(error => {
            alert(error.message);
        });
});

// Обработчик отправки формы регистрации
document.getElementById('signup').addEventListener('submit', function(e) {
    e.preventDefault();

    const formData = new FormData(e.target);

    fetch('/api/signup', {
        method: 'POST',
        body: formData
    })
        .then(response => {
            if (response.ok) {
                return response.json();
            } else {
                return response.json().then(errorData => {
                    throw new Error(errorData.message || 'Registration failed');
                });
            }
        })
        .then(data => {
            alert(data.message);
            showLoginForm();
        })
        .catch(error => {
            alert(error.message);
        });
});

// Обработчики для переключения между формами регистрации и входа
document.getElementById('show-login-link').addEventListener('click', function(e) {
    e.preventDefault();
    showLoginForm();
});

document.getElementById('show-signup-link').addEventListener('click', function(e) {
    e.preventDefault();
    showSignupForm();
});

// Инициализация приложения
function initApp() {
    // Проверяем, авторизован ли пользователь
    fetch('/api/check-auth', {
        method: 'GET',
        credentials: 'include' // Важно для отправки куки вместе с запросом
    })
        .then(response => {
            if (response.ok) {
                showAppSection();
            } else {
                showLoginForm();
            }
        })
        .catch(error => {
            console.error('Error:', error);
            showLoginForm();
        });
}

// Запускаем инициализацию при загрузке страницы
window.addEventListener('load', initApp);
