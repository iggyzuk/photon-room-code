async function makeAjaxCall() {
    var response = await fetch('/api/fruits')
    var text = await response.text();
    alert(`we called out to the api via ajax and got this response => ${text}`);
}