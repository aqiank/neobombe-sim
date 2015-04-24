var scene = new THREE.Scene();
var camera = new THREE.PerspectiveCamera(75, window.innerWidth / window.innerHeight, 0.1, 1000);
var renderer = new THREE.WebGLRenderer();
var ctx = null;
var initialized = false;

var units = [];
var unitFrames = [];
var text = null;

function initUnits(data) {
	if (initialized) {
		return;
	}
	for (var i = 0; i < data.units.length; i++) {
		addUnit(data.units[i], i);
	}
	initialized = true;
}

function updateUnits(data) {
	if (!initialized) {
		return;
	}
	for (var i = 0; i < data.units.length; i++) {
		updateUnit(data.units[i], i);
	}

	drawText(data.text);
}

function update() {
	$.get("/update", function(resp) {
		if (!initialized) {
			initUnits(resp);
		} else {
			updateUnits(resp);
		}
		render();
		update();
	}, "json");
}

function render() {
	renderer.render(scene, camera);
}

function addUnit(unit, idx) {
	var colors = [0x00ff00, 0xffff00, 0xff0000];
	var group = Math.floor(idx / 12);
	var x = idx % 12 - 6;
	var y = -Math.floor(idx / 12) - (group * 1) + 3;
	var rotors = [];
	for (var i = 0; i < unit.rotors.length; i++) {
		var geometry = new THREE.BoxGeometry(0.5, 0.05, 0.01);
		var material = new THREE.MeshBasicMaterial({color: colors[group]});
		var rotor = new THREE.Mesh(geometry, material);
		var xx = (idx % 12) * 0.05;
		var yy = -i * 0.5;
		rotor.position.set(x + xx, y + yy, 0);
		rotors.push(rotor);
		addUnitFrame(idx, x + xx, y + yy);
		scene.add(rotor);
	}
	units.push({rotors: rotors});
}

function addUnitFrame(idx, x, y) {
	var colors = [0x004000, 0x404000, 0x400000];
	var group = Math.floor(idx / 12);
	var geometry = new THREE.TorusGeometry(0.25, 0.01, 32, 24);
	var material = new THREE.MeshBasicMaterial({color: colors[group]});
	var frame = new THREE.Mesh(geometry, material);
	frame.position.set(x, y, 0);
	unitFrames.push(frame);
	scene.add(frame);
}

function updateUnit(unit, idx) {
	for (var i = 0; i < unit.rotors.length; i++) {
		units[idx].rotors[i].rotation.z = unit.rotors[i].angle;
	}
}

function createOverlay() {
	var overlay = document.createElement("canvas");
	overlay.width = window.innerWidth;
	overlay.height = window.innerHeight;
	overlay.style.zIndex = 100;
	overlay.style.position = "absolute";
	document.body.appendChild(overlay);
	ctx = overlay.getContext("2d");
	ctx.fillStyle = "white";
	ctx.font = "16px sans-serif";
}

function drawText(s) {
	ctx.clearRect(0, 0, ctx.canvas.width, ctx.canvas.height);
	ctx.fillText(s, 0, 16);
}

$(document).ready(function() {
	renderer.setSize(window.innerWidth, window.innerHeight);
	document.body.appendChild(renderer.domElement);

	createOverlay();

	camera.position.z = 5;

	update();
});
