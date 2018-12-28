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
  const geometry = new THREE.BoxGeometry(20, 20, 20);
  const material = new THREE.MeshNormalMaterial();
  const box = new THREE.Mesh(geometry, material);
  box.position.set(10, 50, 10)
  scene.add(box);

  tick();

  // 毎フレーム時に実行されるループイベントです
  function tick() {
    // box.rotation.y += 0.01;
    controls.update();
    renderer.render(scene, camera); // レンダリング
    
    requestAnimationFrame(tick);
  }

  setInterval(() =>{
    getPostion().then(position =>{
      position = scaleMotion(position)
      box.position.set(position.x, position.y, position.z)
      box.rotation.set(position.angleX,position.angleY,position.angleZ)
    }).catch((error) => { 
      console.error(`エラーが発生しました (${error})`);
  });
  }, 100)


  var scaleMotion = (position) => {
    var scale = 1000
    return {
      x: position.x * scale,
      y: position.y * scale,
      z: 0,//position.z / scale,
      angleX: position.angleX,
      angleY: position.angleY,
      angleZ: position.angleZ,
    }
  }

  var getPostion = () => {
    return new Promise((resolve, reject) =>{
      var xhttp = new XMLHttpRequest();
      xhttp.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200) {
          // document.getElementById("demo").innerHTML = this.responseText;
          const obj = JSON.parse(this.responseText)
          resolve(obj)
        }
      };
      xhttp.open("GET", "position", true);
      xhttp.send();
    })
  }
}