const search = document.getElementById("inp");
const btn = document.getElementById("btn");

function sendFetch(value) {
    const packag = {
        action: "searchGo",
        data: { value: value }
    };

    fetch("http://127.0.0.1:8080/back/main", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(packag)
    })
    .then(response => {
        if (!response.ok)
            throw new Error("HTTP Error " + response.status + " " + response.statusText);
        return response.json();
    })
    .then(data => { 
        console.log("Réponse du serveur :", data); 
        sessionStorage.setItem(
            "searchData", 
            JSON.stringify(data)
        );
        window.location.href = "inSearch.html"; 
    })
    .catch(error => {
        console.log(`${value} -> not in bdd`);
    });
}

search.addEventListener("keydown", (event) => {
    if (event.key === "Enter") {
        const value = search.value.trim();
        if (value) sendFetch(value);
    }
});

btn.addEventListener("click", () => {
    const value = search.value.trim();
    if (value) sendFetch(value);
});
