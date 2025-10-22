async function getIp() {
  try {
    const res = await fetch("https://api.ipify.org?format=json"); 
    const data = await res.json();
    return data.ip;
  } catch (err) {
    console.error("Erreur:", err);
  }

  return "0.0.0.0";
}

function genUUID (){return crypto.randomUUID();}

export function genOrCreatUUID() {
    const match = document.cookie.match(/(?:^| )userID=([^;]+)/);

    if (match) {
        return match[1]; 
    } else {
        const uuid = genUUID();

        document.cookie = `userID=${uuid}; path=/; max-age=${60 * 60 * 24 * 365 * 100}; SameSite=Lax; Secure`;

        return uuid;
    }
}

window.onload = async () => {
    const toggle = document.getElementById("toggleMode");
    let theme = toggle.checked ? "dark" : "light";

    const ip = await getIp()

    let packag = {
        "action":"chekUsers",
        "data":{
            "ip":String(ip),
            "uuid":String(genOrCreatUUID()),
            "theme":theme
        }
    }

    fetch('http://127.0.0.1:8080/back/main',{
        method:"POST",
        headers:{
            "Content-type":"application/json"
        },
        body:JSON.stringify(packag)
    })

    .then(response => {
        if (!response.ok) throw new Error(
            "HTTP error " + response.status + " " + response.statusText
        );
        return response.json();
    })
    .then(data => {
        console.log(data);
        const body = document.body;
        const modeLabel = document.getElementById('modeLabel');
        const toggle = document.getElementById('toggleMode');

        if (data.status === "old") {
            const theme = data.output === "light" ? "light" : "dark";
            const isDark = theme === "dark";
            
            body.classList.toggle("dark", isDark);
            body.classList.toggle("light", !isDark);
            modeLabel.textContent = isDark ? "Dark Mode" : "Light Mode";
            toggle.checked = isDark;
            sessionStorage.setItem("theme", theme);
        }

        toggle.addEventListener('change', () => {
            const newTheme = toggle.checked ? "dark" : "light";
            body.className = newTheme;
            modeLabel.textContent = toggle.checked ? "Dark Mode" : "Light Mode";
            sessionStorage.setItem("theme", newTheme);
        });
    })
    .catch (error => {
        console.log(error);
    })
}

