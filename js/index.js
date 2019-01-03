// ページの読み込みを待つ
window.addEventListener('load', init);

function init() {

  // サイズを指定
  const width = 960;
  const height = 540;

  // レンダラーを作成
  const renderer = new THREE.WebGLRenderer({
    canvas: document.querySelector('#myCanvas')
  });
  renderer.setPixelRatio(window.devicePixelRatio);
  renderer.setSize(width, height);

  // シーンを作成
  const scene = new THREE.Scene();

  // カメラを作成
  const camera = new THREE.PerspectiveCamera(45, width / height);
//   camera.position.set(0, 0, +1000);
  camera.position.set(100, 150, 500);
  camera.lookAt(new THREE.Vector3(0, 0, 0));
  // camera.up.set(0,0,1);

    // Three.jsのOrbitControl.jsを使ってマウスやタッチで簡単にカメラ操作
    var controls = new THREE.OrbitControls(camera);


  // 地面を作成
  const plane2 = new THREE.GridHelper(600);
  scene.add(plane2);
  const plane = new THREE.AxesHelper(300);
  scene.add(plane);
  const group = new THREE.Group();
  scene.add(group);

  // 箱を作成
  const geometry = new THREE.BoxGeometry(100, 5, 100);
  const material = new THREE.MeshNormalMaterial();
  const box = new THREE.Mesh(geometry, material);
  box.position.set(0, 0, 0)
  scene.add(box);

  tick();

  // 毎フレーム時に実行されるループイベントです
  function tick() {
    // box.rotation.y += 0.01;
    controls.update();
    renderer.render(scene, camera); // レンダリング
    
    requestAnimationFrame(tick);
  }

  var openWebSocket = () => {
    var url = location.hostname + (location.port ? ':' + location.port : '');
    var socket = new WebSocket(`ws://${url}/ws`);
    socket.onopen = function () {
      console.log("onopen")
    };
    socket.onmessage = function (e) {
      updatePosition(JSON.parse(e.data))
    };  
  }

  openWebSocket()

  var updatePosition = (position) => {
    pos = convertAxis(position)
    pos = scaleMotion(pos)
    box.position.set(pos.x, pos.y, pos.z)
    box.rotation.set(pos.angleX,pos.angleY,pos.angleZ)
  }

  var scaleMotion = (position) => {
    var scale = 1/2
    return {
      x: position.x * scale,
      y: position.y * scale,
      z: position.z * scale,
      angleX: position.angleX,
      angleY: position.angleY,
      angleZ: position.angleZ,
    }
  }

  var convertAxis = (position) =>{
    return {
      x: position.x,
      y: position.z,
      z: position.y,
      angleX: position.angleX,
      angleY: position.angleZ,
      angleZ: position.angleY,
    }
  }

  var getPostion = () => {
    return new Promise((resolve, reject) =>{
      var xhttp = new XMLHttpRequest();
      xhttp.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200) {
          const obj = JSON.parse(this.responseText)
          resolve(obj)
        }
      };
      xhttp.open("GET", "position", true);
      xhttp.send();
    })
  }
}