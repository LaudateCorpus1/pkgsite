/*!
 * @license
 * Copyright 2021 The Go Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */const n={PLAY_HREF:".js-exampleHref",PLAY_CONTAINER:".js-exampleContainer",EXAMPLE_INPUT:".Documentation-exampleCode",EXAMPLE_OUTPUT:".Documentation-exampleOutput",EXAMPLE_ERROR:".Documentation-exampleError",PLAY_BUTTON:".Documentation-examplePlayButton",SHARE_BUTTON:".Documentation-exampleShareButton",FORMAT_BUTTON:".Documentation-exampleFormatButton",RUN_BUTTON:".Documentation-exampleRunButton"};export class PlaygroundExampleController{constructor(t){this.exampleEl=t;this.exampleEl=t,this.anchorEl=t.querySelector("a"),this.errorEl=t.querySelector(n.EXAMPLE_ERROR),this.playButtonEl=t.querySelector(n.PLAY_BUTTON),this.shareButtonEl=t.querySelector(n.SHARE_BUTTON),this.formatButtonEl=t.querySelector(n.FORMAT_BUTTON),this.runButtonEl=t.querySelector(n.RUN_BUTTON),this.inputEl=this.makeTextArea(t.querySelector(n.EXAMPLE_INPUT)),this.outputEl=t.querySelector(n.EXAMPLE_OUTPUT),this.playButtonEl?.addEventListener("click",()=>this.handleShareButtonClick()),this.shareButtonEl?.addEventListener("click",()=>this.handleShareButtonClick()),this.formatButtonEl?.addEventListener("click",()=>this.handleFormatButtonClick()),this.runButtonEl?.addEventListener("click",()=>this.handleRunButtonClick()),!!this.inputEl&&(this.resize(),this.inputEl.addEventListener("keyup",()=>this.resize()),this.inputEl.addEventListener("keydown",e=>this.onKeydown(e)))}makeTextArea(t){const e=document.createElement("textarea");return e.classList.add("Documentation-exampleCode"),e.spellcheck=!1,e.value=t?.textContent??"",t?.parentElement?.replaceChild(e,t),e}getAnchorHash(){return this.anchorEl?.hash}expand(){this.exampleEl.open=!0}resize(){if(this.inputEl?.value){const t=(this.inputEl.value.match(/\n/g)||[]).length;this.inputEl.style.height=`${(20+t*20+12+2)/16}rem`}}onKeydown(t){t.key==="Tab"&&(document.execCommand("insertText",!1,"	"),t.preventDefault())}setInputText(t){this.inputEl&&(this.inputEl.value=t)}setOutputText(t){this.outputEl&&(this.outputEl.innerHTML=t)}setErrorText(t){this.errorEl&&(this.errorEl.textContent=t),this.setOutputText("An error has occurred\u2026")}handleShareButtonClick(){const t="https://play.golang.org/p/";this.setOutputText("Waiting for remote server\u2026"),fetch("/play/share",{method:"POST",body:this.inputEl?.value}).then(e=>e.text()).then(e=>{const r=t+e;this.setOutputText(`<a href="${r}">${r}</a>`),window.open(r)}).catch(e=>{this.setErrorText(e)})}handleFormatButtonClick(){this.setOutputText("Waiting for remote server\u2026");const t=new FormData;t.append("body",this.inputEl?.value??""),fetch("/play/fmt",{method:"POST",body:t}).then(e=>e.json()).then(({Body:e,Error:r})=>{this.setOutputText(r||"Done."),e&&(this.setInputText(e),this.resize())}).catch(e=>{this.setErrorText(e)})}handleRunButtonClick(){this.setOutputText("Waiting for remote server\u2026"),fetch("/play/compile",{method:"POST",body:JSON.stringify({body:this.inputEl?.value,version:2})}).then(t=>t.json()).then(async({Events:t,Errors:e})=>{this.setOutputText(e||"");for(const r of t||[])this.setOutputText(r.Message),await new Promise(a=>setTimeout(a,r.Delay/1e6))}).catch(t=>{this.setErrorText(t)})}}const l=location.hash.match(/^#(example-.*)$/);if(l){const o=document.getElementById(l[1]);o&&(o.open=!0)}const i=[...document.querySelectorAll(n.PLAY_HREF)],s=o=>i.find(t=>t.hash===o.getAnchorHash());for(const o of document.querySelectorAll(n.PLAY_CONTAINER)){const t=new PlaygroundExampleController(o),e=s(t);e?e.addEventListener("click",()=>{t.expand()}):console.warn("example href not found")}
//# sourceMappingURL=playground.js.map
