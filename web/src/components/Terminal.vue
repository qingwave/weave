<template>
  <div class="h-full w-full bg" id="terminal"></div>
</template>

<style scoped>
html::-webkit-scrollbar,
body::-webkit-scrollbar,
div::-webkit-scrollbar {
  display: none;
  width: 0;
}

.bg {
  background-image: url("@/assets/code.svg");
  background-repeat: no-repeat;
  background-size: 33%;
  background-position: center;
}
</style>

<script setup>
import "xterm/css/xterm.css";
import { Terminal } from "xterm";
import { AttachAddon } from "xterm-addon-attach";
import { FitAddon } from "xterm-addon-fit";
import { SerializeAddon } from "xterm-addon-serialize";
import { Unicode11Addon } from "xterm-addon-unicode11";
import { WebLinksAddon } from "xterm-addon-web-links";
import { useRoute } from 'vue-router';
import { onMounted, onUnmounted } from 'vue';
import { Windows } from "@icon-park/vue-next";

const props = defineProps({
    type: String,
    uri: String,
    method: String,
})

const route = useRoute();
const term = new Terminal({
  screenKeys: true,
  useStyle: true,
  cursorBlink: true,
  fontFamily: "monospace",
  fullscreenWin: true,
  maximizeWin: true,
  screenReaderMode: true,
  fontSize: 14,
});

const fitAddon = new FitAddon();

let ws = null;

const initTerm = () => {
  const attachAddon = new AttachAddon(ws);
  const serializeAddon = new SerializeAddon();
  const unicode11Addon = new Unicode11Addon();
  const webLinksAddon = new WebLinksAddon();
  term.loadAddon(attachAddon);
  term.loadAddon(fitAddon);
  term.loadAddon(serializeAddon);
  term.loadAddon(unicode11Addon);
  term.loadAddon(webLinksAddon);
  term.open(document.getElementById("terminal"));
  fitAddon.fit();
  term.focus();
}

const initSocket = () => {
  const uri =  `ws://${window.location.hostname}:${window.location.port}${props.uri}`
  ws = new WebSocket(uri);
  ws.onopen = () => {
    initTerm()
  }

  ws.onclose = () => {
    if (term) {
      term.dispose();
    }
    console.log("close socket");
  }

  ws.onerror = (err) => {
    if (term) {
      term.dispose();
    }
    console.log("socket err: ", err);
  }

  // window.onresize = () => {
  //   fitAddon.fit();
  // }
}

onMounted(
  initSocket
)

onUnmounted(
  () => {
    if (!ws) {
      ws.close();
    }

    term.dispose();
  }
)
</script>
