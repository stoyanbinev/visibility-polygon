// http://paulirish.com/2011/requestanimationframe-for-smart-animating/
window.requestAnimFrame = (function() {
  return  window.requestAnimationFrame       || 
    window.webkitRequestAnimationFrame || 
    window.mozRequestAnimationFrame    || 
    window.oRequestAnimationFrame      || 
    window.msRequestAnimationFrame     || 
    function(/* function */ callback, /* DOMElement */ element){
      window.setTimeout(callback, 1000 / 60);
    };
})();

function init() {
  var width = 800;
  var height = 500;
  var observer_x = 250;
  var observer_y = 300;
  var changed = true;
  var polygons = [];
  var segments = [];
  var poly = [];
  var hammer = new Hammer(document.getElementById("canvas"));
  var canvasOffset=$("#canvas").offset();
    var offsetX=canvasOffset.left;
    var offsetY=canvasOffset.top;



    var mouseIsDown=false;
    var dragLight = false;
    var lastX=0;
    var lastY=0;
    var mouseX=0;
    var mouseY=0;


var socket = io();

      var s2 = io("/scene");
      




  hammer.ondrag = function(e) {

    var x = e.position.x;
    var y = e.position.y;
    if (x < 0 || x > width || y < 0 || y > height) return;
    
    observer_x = x/scale;
    observer_y =(height - y)/scale;
    s2.emit('light', JSON.stringify([observer_x,observer_y]), function(data){
          poly = JSON.parse(data);
          update();

        });

    changed = true;
  }; 


} 

  document.getElementById('btnOpen').onclick = function(){
    openFile(function(txt){
        s2.emit('input', txt, function(data){
          console.log(data);
          var s = JSON.parse(data);
           width = s["viewer_x"];
          height = s["viewer_y"];
          scale = 800/Math.max(width, height);
          width = scale*width;
          height = scale*height;
          observer_x = s["light_x"];
          observer_y = s["light_y"];
          poly = s["polygon"];
          polygons = s["segments"];
          segments = convertToSegments(s["segments"]);
          console.log(poly);
          changed = true;
          update();
        });
    });
}
function openFile(callBack){
  var element = document.createElement('input');
  element.setAttribute('type', "file");
  element.setAttribute('id', "btnOpenFile");
  element.onchange = function(){
      readText(this,callBack);
      document.body.removeChild(this);
      }

  element.style.display = 'none';
  document.body.appendChild(element);

  element.click();
}

function readText(filePath,callBack) {
    var reader;
    if (window.File && window.FileReader && window.FileList && window.Blob) {
        reader = new FileReader();
    } else {
        alert('The File APIs are not fully supported by your browser. Fallback required.');
        return false;
    }
    var output = ""; 
    if(filePath.files && filePath.files[0]) {           
        reader.onload = function (e) {
            output = e.target.result;
            callBack(output);
        };
        reader.readAsText(filePath.files[0]);
    }
    else { 
        return false;
    }       
    return true;
}

function distance(aX,aY,bX,bY) {
  var dx = aX - bX;
  var dy = aY - bY;
  return Math.sqrt(dx*dx + dy*dy);
}

  function convertToSegments (polygons) {
  var segments = [];
  for (var i = 0; i < polygons.length; ++i) {
    for (var j = 0; j < polygons[i].length; ++j) {
      var k = j+1;
      if (k == polygons[i].length) k = 0;
      segments.push([[polygons[i][j][0], polygons[i][j][1]], [polygons[i][k][0], polygons[i][k][1]]]);
    }
  }
  return segments;
};

  function update() {
    if (changed) {
      var canvas = document.getElementById('canvas');
      var ctx = canvas.getContext("2d");

            ctx.clearRect(0, 0, width, height);
      ctx.beginPath();
      ctx.rect(0, 0, width, height);
      ctx.fillStyle = '#666';
      ctx.fill();
           draw(ctx);
           changed = false;

      
      
    }
    requestAnimFrame(update);
    
  };

  function draw(ctx) {


    for (var i = 1; i < polygons.length; ++i) {
      ctx.beginPath();
      ctx.moveTo(polygons[i][0][0]*scale, height - polygons[i][0][1]*scale);
      for (var j = 1; j < polygons[i].length; ++j) {
        ctx.lineTo(polygons[i][j][0]*scale, height - polygons[i][j][1]*scale);
      }
      ctx.fillStyle = "orange";
      ctx.fill();
    }

    ctx.beginPath();
    ctx.moveTo(poly[0][0]*scale, height - poly[0][1]*scale);
    for (var i = 1; i < poly.length; ++i) {
      ctx.lineTo(poly[i][0]*scale, height - poly[i][1]*scale);
    }
    ctx.fillStyle = "#aaa";
    ctx.fill();

    for (var i = 0; i < segments.length; ++i) {
      ctx.beginPath();
      ctx.moveTo(segments[i][0][0]*scale, height - segments[i][0][1]*scale);
      ctx.lineTo(segments[i][1][0]*scale, height - segments[i][1][1]*scale);
      ctx.strokeStyle = "black";
      ctx.lineWidth = 2;
      ctx.stroke();
    }

    ctx.beginPath();
    for (var i = 0; i < segments.length; ++i) {
      ctx.beginPath();
    
    ctx.arc(segments[i][0][0]*scale, height - segments[i][0][1]*scale, 5, 0, Math.PI*2, true);
    ctx.fillStyle = "black";
    ctx.fill();
    ctx.strokeStyle = "black";
    ctx.stroke();
  }

    ctx.beginPath();
    ctx.arc(observer_x*scale, height - observer_y*scale, 5, 0, Math.PI*2, true);
    ctx.fillStyle = "yellow";
    ctx.fill();
    ctx.strokeStyle = "black";
    ctx.stroke();
  };



  function handleMouseDown(e){

    console.log("Down");
  // get the current mouse position relative to the canvas

  mouseX=parseInt(e.clientX-offsetX);
  mouseY=parseInt(e.clientY-offsetY);


  if(distance(mouseX,height-mouseY,observer_x*scale,observer_y*scale)<10){
    dragLight = true;
  }

  // save this last mouseX/mouseY

console.log(mouseX,mouseY);
  lastX=mouseX;
  lastY=mouseY;

  // set the mouseIsDown flag

  mouseIsDown=true;
}


function handleMouseUp(e){

console.log("Up");
  // clear the mouseIsDown flag
  if(dragLight){
    console.log()
    observe_x = mouseX;
    observe_y = height-mouseY;
    update();
  }

  dragLight = false;
  mouseIsDown=false;
}

function handleMouseMove(e){

  if(!mouseIsDown){ return; }

  mouseX=parseInt(e.clientX-offsetX);
  mouseY=parseInt(e.clientY-offsetY);

  lastX=mouseX;
  lastY=mouseY;
  
}




};