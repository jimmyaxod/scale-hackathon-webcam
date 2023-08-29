
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

  let frameTimes:number[] = [];
  let maxFrameTimes = 1000;

  let el_fps = document.getElementById("fps");

  const doSomethingWithTheFrame = (now: any, metadata: any) => {
    frameTimes.push(now);
    if (frameTimes.length > maxFrameTimes) {
      frameTimes.unshift();
    }

    let duration = now - frameTimes[0];
    if (duration>0 && el_fps!=null) {
      let fps = frameTimes.length * 1000 / duration;
      el_fps.innerText = "FPS: " + fps.toFixed(2);
    }

    let ctx = canvas.getContext('2d');
    if (ctx!=null) {
      ctx.drawImage(video, 0, 0, canvas.width, canvas.height);
    }

    // Do something with the frame.
    //console.log("Video frame: " + now, metadata);
    // Re-register the callback to be notified about the next frame.
    video.requestVideoFrameCallback(doSomethingWithTheFrame);
  };
  // Initially register the callback to be notified about the first frame.
  video.requestVideoFrameCallback(doSomethingWithTheFrame);

  navigator.mediaDevices.getUserMedia( {audio: false, video: true })
    .then(function(stream) {
	    console.log("Setting video srcObject:", stream);
	    video.srcObject = stream;
    })
    .catch(error => console.error(error));
  }

})
