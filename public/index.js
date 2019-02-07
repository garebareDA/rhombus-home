'use strict';
var so = io();

window.SpeechRecognition = window.webkitSpeechRecognition || window.SpeechRecognition;
const recognition = new SpeechRecognition();
recognition.lang = 'ja-JP';
recognition.continuous = true;

var ss = new SpeechSynthesisUtterance();
ss.lang = "ja-JP";
ss.voiceURI = "Google 日本人";
ss.volume = 1;
ss.rate = 1;
ss.pitch = 1;

// レンダラの作成、DOMに追加
const renderer = new THREE.WebGLRenderer();
const clock = new THREE.Clock();
let mixer;

document.getElementById('child').appendChild(renderer.domElement);

const scene = new THREE.Scene();
const camera = new THREE.PerspectiveCamera(45, 1.0);
camera.position.set(0, 1, 5);

renderer.setClearColor(0xf3f3f3, 1.0);

var directionalLight = new THREE.DirectionalLight( 0xffeedd, 10);
directionalLight.position.set( 0, 0, 50 );
scene.add( directionalLight );

var directionalLight2 = new THREE.DirectionalLight( 0xffeedd, 10);
directionalLight2.position.set( 0, 0, -50 );
scene.add( directionalLight2 );

var directionalLight3 = new THREE.DirectionalLight( 0xffeedd, 5);
directionalLight3.position.set( 50, 0, 0 );
scene.add( directionalLight3 );

var directionalLight4 = new THREE.DirectionalLight( 0xffeedd, 5);
directionalLight4.position.set( -50, 0, 0 );
scene.add( directionalLight4 );

var directionalLight5 = new THREE.DirectionalLight( 0xffeedd, 5);
directionalLight5.position.set( 0, 50, 0 );
scene.add( directionalLight5 );

const loader = new THREE.GLTFLoader();
const model = 'public/assistant.gltf';
window.addEventListener('resize', onResize);

loader.load(model,(data) => {
  const gltf = data;
  const object = gltf.scene;
  const animations = gltf.animations

  mixer = new THREE.AnimationMixer(object);
  object.position.set(0, 1.2, -8,)
  scene.add(object);

  const child = document.getElementById('child');
  const click = document.getElementById('click');
  const oneceButton = function(){
    click.remove();
    helloAnime();
    speak('こんにちは');
    recognition.start();
    child.removeEventListener('click', oneceButton);
  }
  child.addEventListener('click', oneceButton);

  so.on('message',function(msg){
    if(msg != null){
      talkAnime();
      speak(msg);
      ss.onend = () =>{
        const anime = mixer.clipAction(animations[2]);
        anime.setLoop(THREE.LoopOnce);
        anime.clampWhenFinished = true;
        anime.play();
      }

    }else{
      talkAnime();
      speak("すみませんエラーが起きました");
      ss.onend = () =>{
        const anime = mixer.clipAction(animations[2]);
        anime.setLoop(THREE.LoopOnce);
        anime.clampWhenFinished = true;
        anime.play();
      }
    }
  });

  recognition.onresult = (event) => {
    const text = event.results[0][0].transcript;
    so.emit('message', text);
    console.log(text);
    recognition.stop();
    mixer.stopAllAction();
    mixer.clipAction(animations[0]).play();
  }

  function helloAnime() {
    mixer.stopAllAction();
    const anime = mixer.clipAction(animations[1]);
    anime.setLoop(THREE.LoopOnce);
    anime.clampWhenFinished = false;
    anime.play();
  }

  function talkAnime(){
    mixer.stopAllAction();
    const anime = mixer.clipAction(animations[2]);
    anime.setLoop(THREE.LoopRepeat);
    anime.clampWhenFinished = true;
    anime.play();
  }
});

//常に認識
recognition.onend = (event) => {
  recognition.start()
}

animation();
onResize();

//リサイズ
function onResize() {
  const width = window.innerWidth;
  const height = window.innerHeight;

  // レンダラーのサイズを調整する
  renderer.setPixelRatio(window.devicePixelRatio);
  renderer.setSize(width, height);

  // カメラのアスペクト比を正す
  camera.aspect = width / height;
  camera.updateProjectionMatrix();
}

function animation() {
  renderer.render(scene, camera);

  if (mixer) {
    mixer.update(clock.getDelta());
  }

  requestAnimationFrame(animation);
};

function speak(text) {
  ss.text = text
  speechSynthesis.speak(ss);
}