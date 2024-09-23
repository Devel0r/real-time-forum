// Получаем ссылки на элементы
const loginForm = document.getElementById('login-form');
const signupForm = document.getElementById('signup-form');
const postsSection = document.getElementById('posts-section');
const chatSection = document.getElementById('chat-section');
const authSection = document.getElementById('auth-section');
const loginBtn = document.getElementById('login-btn');
const signupBtn = document.getElementById('signup-btn');
const postsBtn = document.getElementById('posts-btn');
const logoutBtn = document.getElementById('logout-btn');

// Функции для показа/скрытия секций
function showLoginForm() {
    loginForm.style.display = 'block';
    signupForm.style.display = 'none';
    postsSection.style.display = 'none';
    chatSection.style.display = 'none';
}

function showSignupForm() {
    signupForm.style.display = 'block';
    loginForm.style.display = 'none';
    postsSection.style.display = 'none';
    chatSection.style.display = 'none';
}

function showPostsSection() {
    postsSection.style.display = 'block';
    loginForm.style.display = 'none';
    signupForm.style.display = 'none';
    authSection.style.display = 'none';
    chatSection.style.display = 'none';
}

function showChatSection() {
    chatSection.style.display = 'block';
    postsSection.style.display = 'none';
    authSection.style.display = 'none';
    loginForm.style.display = 'none';
    signupForm.style.display = 'none';
}

// Слушатели событий для кнопок
loginBtn.addEventListener('click', showLoginForm);
signupBtn.addEventListener('click', showSignupForm);
postsBtn.addEventListener('click', showPostsSection);
logoutBtn.addEventListener('click', () => {
    // Логика выхода
    authSection.style.display = 'block';
    postsSection.style.display = 'none';
    chatSection.style.display = 'none';
    postsBtn.style.display = 'none';
    logoutBtn.style.display = 'none';
});

// Пример входа в систему
document.getElementById('login').addEventListener('submit', function(e) {
    e.preventDefault();
    // Здесь будет логика для авторизации пользователя
    alert('Logged in!');
    authSection.style.display = 'none';
    postsBtn.style.display = 'inline-block';
    logoutBtn.style.display = 'inline-block';
    showPostsSection();
});

// Пример регистрации
document.getElementById('signup').addEventListener('submit', function(e) {
    e.preventDefault();
    // Логика для регистрации пользователя
    alert('Signed up!');
    showLoginForm();
});