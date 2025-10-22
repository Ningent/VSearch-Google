const inp = document.getElementById("inp")
const btnn = document.getElementById("btn")

function sendFetch (value) {
    console.log(1)
    let packag = {
        "action":"newSearch",
        "data":{
            "search":value
        }
    }

    fetch ("http://127.0.0.1:8080/back/main",{
        method:"POST",
        headers:{
            "Content-Type":"application/json"
        },

        body:JSON.stringify(packag)
    })
    .then (response => {
        if (!response.ok)
            throw new Error(
                `HTTP ERROR 
                ${response.status}\n
                ${response.statusText}
            `)

        return response.json()
    })

    .then (data => {
        console.log (`new Result ${data}`)
    })

    .catch (error => {
        console.error(error)
    })

}


inp.addEventListener("keydown",(event) => {
    if (event.key == "Enter"){
        const value = inp.value.trim()
        if  (value) sendFetch(value);
    }
})


btnn.addEventListener ("click",(event) => {
    const value = inp.value.trim()
    if (value) sendFetch(value)
})