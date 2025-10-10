import { genOrCreatUUID } from './userAuthen.js';

const toggle = document.getElementById('toggleMode');
const body = document.body;
const modeLabel = document.getElementById('modeLabel');

body.classList.add('light');
let theme = "light";



toggle.addEventListener('change', async () => {
    if (toggle.checked) {
        body.classList.remove('light');
        body.classList.add('dark');
        modeLabel.textContent = "Dark Mode";
        theme = "dark";
    } else {
        body.classList.remove('dark');
        body.classList.add('light');
        modeLabel.textContent = "Light Mode";
        theme = "light";
    }
  
    

    const packag = {
        action: "changeTheme",
        data: {
            uuid: String(genOrCreatUUID()),
            theme: theme
        }
    };


    fetch("http://127.0.0.1:8080/back/main", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(packag)
    })
    .then(response => {
        if (!response.ok) throw new Error(
            "HTTP error " + response.status + " " + response.statusText
        );
        return response.json();
    })
    .then(data => {
        console.log("reponse du server", data);
    })
    .catch(error => {
        console.error("Erreur fetch:", error);
    });
});