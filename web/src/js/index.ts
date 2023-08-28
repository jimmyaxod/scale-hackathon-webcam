
document.addEventListener('DOMContentLoaded', async function() {

  const video = document.getElementById('video') as HTMLVideoElement;
  const canvas = document.getElementById('canvas') as HTMLCanvasElement;
  const button = document.getElementById('button') as HTMLButtonElement;
  if (canvas!=null && video!=null && button!=null) {
    canvas.width = 480;
    canvas.height = 360;


  button.onclick = function() {
  	console.log("Getting snapshot");
    canvas.width = video.videoWidth;
    canvas.height = video.videoHeight;
    let ctx = canvas.getContext('2d');
    if (ctx!=null) {
      ctx.drawImage(video, 0, 0, canvas.width, canvas.height);
    }
  };

  navigator.mediaDevices.getUserMedia( {audio: false, video: true })
    .then(function(stream) {
	    console.log("Setting video srcObject:", stream);
	    video.srcObject = stream;
    })
    .catch(error => console.error(error));
  }

})
