const div = document.getElementById("popup");
const settings = document.getElementsByClassName("settings")[0];
settings.onclick = function(event) {
    if (div.style.display != "block")
        div.style = "display: block";
    else {
        div.style = "display: none";
    }
};

document.addEventListener( 'click', (e) => {
	const withinBoundaries_1 = e.composedPath().includes(settings);
    const withinBoundaries_2 = e.composedPath().includes(div);
	if (! withinBoundaries_1 && ! withinBoundaries_2) {
		div.style.display = 'none';
		
	}
})

document.addEventListener('keydown', function(e) {
	if(e.key === "Escape"){
		div.style.display = 'none';
	}
});