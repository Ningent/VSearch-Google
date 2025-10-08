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

function genOrCreatUUID(){
    let match = document.cookie.match(/(^| )userID=([^;]+)/)

    if (match){
        return match[2];
    }else {
        const uuid = genUUID();

        document.cookie = `
            userID=${uuid}; 
            path=/; 
            max-age=${60 * 60 * 24 * 365*100}
        `;

        return uuid;
    }   
}

window.onload = async () => {
    const toggle = document.getElementById("toggleMode");
    let theme = toggle.checked ? "dark" : "light";

    const ip = await getIp()

    packag = {
        "action":"NewUser",
        "data":{
            "ip":ip,
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

    .then (response => response.json())
    .then (data => {
        console.log (data);
    })
    .catch (error => {
        console.log(error);
    })
}
