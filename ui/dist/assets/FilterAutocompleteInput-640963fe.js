import{S as oe,i as ae,s as le,e as ue,f as ce,g as fe,y as M,o as de,I as he,J as ge,K as pe,L as ye,C as S,M as me}from"./index-08378882.js";import{C as E,E as q,a as w,h as be,b as ke,c as xe,d as Ke,e as Ce,s as Se,f as qe,g as we,i as Le,r as Ie,j as Ee,k as Re,l as Ae,m as Be,n as Oe,o as _e,p as ve,q as Me,t as Y,S as De}from"./index-a6ccb683.js";function He(e){Z(e,"start");var i={},n=e.languageData||{},g=!1;for(var p in e)if(p!=n&&e.hasOwnProperty(p))for(var d=i[p]=[],o=e[p],s=0;s<o.length;s++){var l=o[s];d.push(new Ue(l,e)),(l.indent||l.dedent)&&(g=!0)}return{name:n.name,startState:function(){return{state:"start",pending:null,indent:g?[]:null}},copyState:function(a){var b={state:a.state,pending:a.pending,indent:a.indent&&a.indent.slice(0)};return a.stack&&(b.stack=a.stack.slice(0)),b},token:We(i),indent:Ne(i,n),languageData:n}}function Z(e,i){if(!e.hasOwnProperty(i))throw new Error("Undefined state "+i+" in simple mode")}function Fe(e,i){if(!e)return/(?:)/;var n="";return e instanceof RegExp?(e.ignoreCase&&(n="i"),e=e.source):e=String(e),new RegExp((i===!1?"":"^")+"(?:"+e+")",n)}function Te(e){if(!e)return null;if(e.apply)return e;if(typeof e=="string")return e.replace(/\./g," ");for(var i=[],n=0;n<e.length;n++)i.push(e[n]&&e[n].replace(/\./g," "));return i}function Ue(e,i){(e.next||e.push)&&Z(i,e.next||e.push),this.regex=Fe(e.regex),this.token=Te(e.token),this.data=e}function We(e){return function(i,n){if(n.pending){var g=n.pending.shift();return n.pending.length==0&&(n.pending=null),i.pos+=g.text.length,g.token}for(var p=e[n.state],d=0;d<p.length;d++){var o=p[d],s=(!o.data.sol||i.sol())&&i.match(o.regex);if(s){o.data.next?n.state=o.data.next:o.data.push?((n.stack||(n.stack=[])).push(n.state),n.state=o.data.push):o.data.pop&&n.stack&&n.stack.length&&(n.state=n.stack.pop()),o.data.indent&&n.indent.push(i.indentation()+i.indentUnit),o.data.dedent&&n.indent.pop();var l=o.token;if(l&&l.apply&&(l=l(s)),s.length>2&&o.token&&typeof o.token!="string"){n.pending=[];for(var a=2;a<s.length;a++)s[a]&&n.pending.push({text:s[a],token:o.token[a-1]});return i.backUp(s[0].length-(s[1]?s[1].length:0)),l[0]}else return l&&l.join?l[0]:l}}return i.next(),null}}function Ne(e,i){return function(n,g){if(n.indent==null||i.dontIndentStates&&i.doneIndentState.indexOf(n.state)>-1)return null;var p=n.indent.length-1,d=e[n.state];e:for(;;){for(var o=0;o<d.length;o++){var s=d[o];if(s.data.dedent&&s.data.dedentIfLineStart!==!1){var l=s.regex.exec(g);if(l&&l[0]){p--,(s.next||s.push)&&(d=e[s.next||s.push]),g=g.slice(l[0].length);continue e}}}break}return p<0?0:n.indent[p]}}function Je(e){let i;return{c(){i=ue("div"),ce(i,"class","code-editor")},m(n,g){fe(n,i,g),e[15](i)},p:M,i:M,o:M,d(n){n&&de(i),e[15](null)}}}function Pe(e){return JSON.stringify([e==null?void 0:e.name,e==null?void 0:e.type,e==null?void 0:e.schema])}function Ve(e,i,n){let g;he(e,ge,t=>n(21,g=t));const p=pe();let{id:d=""}=i,{value:o=""}=i,{disabled:s=!1}=i,{placeholder:l=""}=i,{baseCollection:a=null}=i,{singleLine:b=!1}=i,{extraAutocompleteKeys:R=[]}=i,{disableRequestKeys:x=!1}=i,{disableIndirectCollectionsKeys:K=!1}=i,f,k,A=s,D=new E,H=new E,F=new E,T=new E,L=[],U=[],W=[],N=[],I="",B="";function O(){f==null||f.focus()}let _=null;function j(){clearTimeout(_),_=setTimeout(()=>{L=$(g),N=ee(),U=x?[]:te(),W=K?[]:ne()},300)}function $(t){let r=t.slice();return a&&S.pushOrReplaceByKey(r,a,"id"),r}function J(){k==null||k.dispatchEvent(new CustomEvent("change",{detail:{value:o},bubbles:!0}))}function P(){if(!d)return;const t=document.querySelectorAll('[for="'+d+'"]');for(let r of t)r.removeEventListener("click",O)}function V(){if(!d)return;P();const t=document.querySelectorAll('[for="'+d+'"]');for(let r of t)r.addEventListener("click",O)}function C(t,r="",c=0){var m,z,Q;let h=L.find(y=>y.name==t||y.id==t);if(!h||c>=4)return[];let u=S.getAllCollectionIdentifiers(h,r);for(const y of h.schema){const v=r+y.name;if(y.type==="relation"&&((m=y.options)!=null&&m.collectionId)){const X=C(y.options.collectionId,v+".",c+1);X.length&&(u=u.concat(X))}y.type==="select"&&((z=y.options)==null?void 0:z.maxSelect)!=1&&u.push(v+":each"),((Q=y.options)==null?void 0:Q.maxSelect)!=1&&["select","file","relation"].includes(y.type)&&u.push(v+":length")}return u}function ee(){return C(a==null?void 0:a.name)}function te(){const t=[];t.push("@request.method"),t.push("@request.query."),t.push("@request.data."),t.push("@request.auth.id"),t.push("@request.auth.collectionId"),t.push("@request.auth.collectionName"),t.push("@request.auth.verified"),t.push("@request.auth.username"),t.push("@request.auth.email"),t.push("@request.auth.emailVisibility"),t.push("@request.auth.created"),t.push("@request.auth.updated");const r=L.filter(h=>h.isAuth);for(const h of r){const u=C(h.id,"@request.auth.");for(const m of u)S.pushUnique(t,m)}const c=["created","updated"];if(a!=null&&a.id){const h=C(a.name,"@request.data.");for(const u of h){t.push(u);const m=u.split(".");m.length===3&&m[2].indexOf(":")===-1&&!c.includes(m[2])&&t.push(u+":isset")}}return t}function ne(){const t=[];for(const r of L){const c="@collection."+r.name+".",h=C(r.name,c);for(const u of h)t.push(u)}return t}function ie(t=!0,r=!0){let c=[].concat(R);return c=c.concat(N||[]),t&&(c=c.concat(U||[])),r&&(c=c.concat(W||[])),c.sort(function(h,u){return u.length-h.length}),c}function se(t){let r=t.matchBefore(/[\'\"\@\w\.]*/);if(r&&r.from==r.to&&!t.explicit)return null;let c=[{label:"false"},{label:"true"},{label:"@now"}];K||c.push({label:"@collection.*",apply:"@collection."});const h=ie(!x,!x&&r.text.startsWith("@c"));for(const u of h)c.push({label:u.endsWith(".")?u+"*":u,apply:u});return{from:r.from,options:c}}function G(){return De.define(He({start:[{regex:/true|false|null/,token:"atom"},{regex:/"(?:[^\\]|\\.)*?(?:"|$)/,token:"string"},{regex:/'(?:[^\\]|\\.)*?(?:'|$)/,token:"string"},{regex:/0x[a-f\d]+|[-+]?(?:\.\d+|\d+\.?\d*)(?:e[-+]?\d+)?/i,token:"number"},{regex:/\&\&|\|\||\=|\!\=|\~|\!\~|\>|\<|\>\=|\<\=/,token:"operator"},{regex:/[\{\[\(]/,indent:!0},{regex:/[\}\]\)]/,dedent:!0},{regex:/\w+[\w\.]*\w+/,token:"keyword"},{regex:S.escapeRegExp("@now"),token:"keyword"},{regex:S.escapeRegExp("@request.method"),token:"keyword"}]}))}ye(()=>{const t={key:"Enter",run:r=>{b&&p("submit",o)}};return V(),n(11,f=new q({parent:k,state:w.create({doc:o,extensions:[be(),ke(),xe(),Ke(),Ce(),w.allowMultipleSelections.of(!0),Se(qe,{fallback:!0}),we(),Le(),Ie(),Ee(),Re.of([t,...Ae,...Be,Oe.find(r=>r.key==="Mod-d"),..._e,...ve]),q.lineWrapping,Me({override:[se],icons:!1}),T.of(Y(l)),H.of(q.editable.of(!s)),F.of(w.readOnly.of(s)),D.of(G()),w.transactionFilter.of(r=>b&&r.newDoc.lines>1?[]:r),q.updateListener.of(r=>{!r.docChanged||s||(n(1,o=r.state.doc.toString()),J())})]})})),()=>{clearTimeout(_),P(),f==null||f.destroy()}});function re(t){me[t?"unshift":"push"](()=>{k=t,n(0,k)})}return e.$$set=t=>{"id"in t&&n(2,d=t.id),"value"in t&&n(1,o=t.value),"disabled"in t&&n(3,s=t.disabled),"placeholder"in t&&n(4,l=t.placeholder),"baseCollection"in t&&n(5,a=t.baseCollection),"singleLine"in t&&n(6,b=t.singleLine),"extraAutocompleteKeys"in t&&n(7,R=t.extraAutocompleteKeys),"disableRequestKeys"in t&&n(8,x=t.disableRequestKeys),"disableIndirectCollectionsKeys"in t&&n(9,K=t.disableIndirectCollectionsKeys)},e.$$.update=()=>{e.$$.dirty[0]&32&&n(13,I=Pe(a)),e.$$.dirty[0]&25352&&!s&&(B!=I||x!==-1||K!==-1)&&(n(14,B=I),j()),e.$$.dirty[0]&4&&d&&V(),e.$$.dirty[0]&2080&&f&&a!=null&&a.schema&&f.dispatch({effects:[D.reconfigure(G())]}),e.$$.dirty[0]&6152&&f&&A!=s&&(f.dispatch({effects:[H.reconfigure(q.editable.of(!s)),F.reconfigure(w.readOnly.of(s))]}),n(12,A=s),J()),e.$$.dirty[0]&2050&&f&&o!=f.state.doc.toString()&&f.dispatch({changes:{from:0,to:f.state.doc.length,insert:o}}),e.$$.dirty[0]&2064&&f&&typeof l<"u"&&f.dispatch({effects:[T.reconfigure(Y(l))]})},[k,o,d,s,l,a,b,R,x,K,O,f,A,I,B,re]}class Qe extends oe{constructor(i){super(),ae(this,i,Ve,Je,le,{id:2,value:1,disabled:3,placeholder:4,baseCollection:5,singleLine:6,extraAutocompleteKeys:7,disableRequestKeys:8,disableIndirectCollectionsKeys:9,focus:10},null,[-1,-1])}get focus(){return this.$$.ctx[10]}}export{Qe as default};
