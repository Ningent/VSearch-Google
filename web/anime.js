const toggle = document.getElementById('toggleMode');
const body = document.body;
const modeLabel = document.getElementById('modeLabel');

body.classList.add('light');

toggle.addEventListener('change', () => {
    if(toggle.checked) {
        body.classList.remove('light');
        body.classList.add('dark');
        modeLabel.textContent = "Dark Mode";
    } else {
        body.classList.remove('dark');
        body.classList.add('light');
        modeLabel.textContent = "Light Mode";
    }
});