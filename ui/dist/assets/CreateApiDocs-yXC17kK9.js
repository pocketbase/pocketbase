import{S as $t,i as qt,s as Tt,V as St,X as ce,W as Ct,h as o,d as $e,t as he,a as ve,I as ae,Z as Ne,_ as pt,C as Mt,$ as Pt,D as Lt,l as r,n as i,m as qe,u as a,A as b,v as p,c as Te,w,J as we,p as Ft,k as Se,o as Ht,L as Ot,H as fe}from"./index-CkwOC79g.js";import{F as Rt}from"./FieldsQueryParam-haBHHZQ1.js";function mt(s,e,t){const l=s.slice();return l[10]=e[t],l}function bt(s,e,t){const l=s.slice();return l[10]=e[t],l}function _t(s,e,t){const l=s.slice();return l[15]=e[t],l}function kt(s){let e;return{c(){e=a("p"),e.innerHTML="Requires superuser <code>Authorization:TOKEN</code> header",w(e,"class","txt-hint txt-sm txt-right")},m(t,l){r(t,e,l)},d(t){t&&o(e)}}}function yt(s){let e,t,l,c,f,u,_,m,q,y,g,B,S,$,R,P,I,D,M,W,L,T,k,F,ee,z,U,oe,K,X,Y;function ue(h,C){var N,x,O;return C&1&&(u=null),u==null&&(u=!!((O=(x=(N=h[0])==null?void 0:N.fields)==null?void 0:x.find(Yt))!=null&&O.required)),u?Bt:At}let te=ue(s,-1),E=te(s);function Z(h,C){var N,x,O;return C&1&&(I=null),I==null&&(I=!!((O=(x=(N=h[0])==null?void 0:N.fields)==null?void 0:x.find(Xt))!=null&&O.required)),I?Nt:Vt}let G=Z(s,-1),H=G(s);return{c(){e=a("tr"),e.innerHTML='<td colspan="3" class="txt-hint txt-bold">Auth specific fields</td>',t=p(),l=a("tr"),c=a("td"),f=a("div"),E.c(),_=p(),m=a("span"),m.textContent="email",q=p(),y=a("td"),y.innerHTML='<span class="label">String</span>',g=p(),B=a("td"),B.textContent="Auth record email address.",S=p(),$=a("tr"),R=a("td"),P=a("div"),H.c(),D=p(),M=a("span"),M.textContent="emailVisibility",W=p(),L=a("td"),L.innerHTML='<span class="label">Boolean</span>',T=p(),k=a("td"),k.textContent="Whether to show/hide the auth record email when fetching the record data.",F=p(),ee=a("tr"),ee.innerHTML='<td><div class="inline-flex"><span class="label label-success">Required</span> <span>password</span></div></td> <td><span class="label">String</span></td> <td>Auth record password.</td>',z=p(),U=a("tr"),U.innerHTML='<td><div class="inline-flex"><span class="label label-success">Required</span> <span>passwordConfirm</span></div></td> <td><span class="label">String</span></td> <td>Auth record password confirmation.</td>',oe=p(),K=a("tr"),K.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>verified</span></div></td> <td><span class="label">Boolean</span></td> <td>Indicates whether the auth record is verified or not.
                    <br/>
                    This field can be set only by superusers or auth records with &quot;Manage&quot; access.</td>`,X=p(),Y=a("tr"),Y.innerHTML='<td colspan="3" class="txt-hint txt-bold">Other fields</td>',w(f,"class","inline-flex"),w(P,"class","inline-flex")},m(h,C){r(h,e,C),r(h,t,C),r(h,l,C),i(l,c),i(c,f),E.m(f,null),i(f,_),i(f,m),i(l,q),i(l,y),i(l,g),i(l,B),r(h,S,C),r(h,$,C),i($,R),i(R,P),H.m(P,null),i(P,D),i(P,M),i($,W),i($,L),i($,T),i($,k),r(h,F,C),r(h,ee,C),r(h,z,C),r(h,U,C),r(h,oe,C),r(h,K,C),r(h,X,C),r(h,Y,C)},p(h,C){te!==(te=ue(h,C))&&(E.d(1),E=te(h),E&&(E.c(),E.m(f,_))),G!==(G=Z(h,C))&&(H.d(1),H=G(h),H&&(H.c(),H.m(P,D)))},d(h){h&&(o(e),o(t),o(l),o(S),o($),o(F),o(ee),o(z),o(U),o(oe),o(K),o(X),o(Y)),E.d(),H.d()}}}function At(s){let e;return{c(){e=a("span"),e.textContent="Optional",w(e,"class","label label-warning")},m(t,l){r(t,e,l)},d(t){t&&o(e)}}}function Bt(s){let e;return{c(){e=a("span"),e.textContent="Required",w(e,"class","label label-success")},m(t,l){r(t,e,l)},d(t){t&&o(e)}}}function Vt(s){let e;return{c(){e=a("span"),e.textContent="Optional",w(e,"class","label label-warning")},m(t,l){r(t,e,l)},d(t){t&&o(e)}}}function Nt(s){let e;return{c(){e=a("span"),e.textContent="Required",w(e,"class","label label-success")},m(t,l){r(t,e,l)},d(t){t&&o(e)}}}function jt(s){let e;return{c(){e=a("span"),e.textContent="Required",w(e,"class","label label-success")},m(t,l){r(t,e,l)},d(t){t&&o(e)}}}function Jt(s){let e;return{c(){e=a("span"),e.textContent="Optional",w(e,"class","label label-warning")},m(t,l){r(t,e,l)},d(t){t&&o(e)}}}function Dt(s){let e,t=s[15].maxSelect===1?"id":"ids",l,c;return{c(){e=b("Relation record "),l=b(t),c=b(".")},m(f,u){r(f,e,u),r(f,l,u),r(f,c,u)},p(f,u){u&32&&t!==(t=f[15].maxSelect===1?"id":"ids")&&ae(l,t)},d(f){f&&(o(e),o(l),o(c))}}}function Et(s){let e,t,l,c,f,u,_,m,q;return{c(){e=b("File object."),t=a("br"),l=b(`
                        Set to empty value (`),c=a("code"),c.textContent="null",f=b(", "),u=a("code"),u.textContent='""',_=b(" or "),m=a("code"),m.textContent="[]",q=b(`) to delete
                        already uploaded file(s).`)},m(y,g){r(y,e,g),r(y,t,g),r(y,l,g),r(y,c,g),r(y,f,g),r(y,u,g),r(y,_,g),r(y,m,g),r(y,q,g)},p:fe,d(y){y&&(o(e),o(t),o(l),o(c),o(f),o(u),o(_),o(m),o(q))}}}function It(s){let e,t;return{c(){e=a("code"),e.textContent='{"lon":x,"lat":y}',t=b(" object.")},m(l,c){r(l,e,c),r(l,t,c)},p:fe,d(l){l&&(o(e),o(t))}}}function Ut(s){let e;return{c(){e=b("URL address.")},m(t,l){r(t,e,l)},p:fe,d(t){t&&o(e)}}}function Qt(s){let e;return{c(){e=b("Email address.")},m(t,l){r(t,e,l)},p:fe,d(t){t&&o(e)}}}function Wt(s){let e;return{c(){e=b("JSON array or object.")},m(t,l){r(t,e,l)},p:fe,d(t){t&&o(e)}}}function xt(s){let e;return{c(){e=b("Number value.")},m(t,l){r(t,e,l)},p:fe,d(t){t&&o(e)}}}function zt(s){let e,t,l=s[15].autogeneratePattern&&ht();return{c(){e=b(`Plain text value.
                        `),l&&l.c(),t=Ot()},m(c,f){r(c,e,f),l&&l.m(c,f),r(c,t,f)},p(c,f){c[15].autogeneratePattern?l||(l=ht(),l.c(),l.m(t.parentNode,t)):l&&(l.d(1),l=null)},d(c){c&&(o(e),o(t)),l&&l.d(c)}}}function ht(s){let e;return{c(){e=b("It is autogenerated if not set.")},m(t,l){r(t,e,l)},d(t){t&&o(e)}}}function vt(s,e){let t,l,c,f,u,_=e[15].name+"",m,q,y,g,B=we.getFieldValueType(e[15])+"",S,$,R,P;function I(k,F){return!k[15].required||k[15].type=="text"&&k[15].autogeneratePattern?Jt:jt}let D=I(e),M=D(e);function W(k,F){if(k[15].type==="text")return zt;if(k[15].type==="number")return xt;if(k[15].type==="json")return Wt;if(k[15].type==="email")return Qt;if(k[15].type==="url")return Ut;if(k[15].type==="geoPoint")return It;if(k[15].type==="file")return Et;if(k[15].type==="relation")return Dt}let L=W(e),T=L&&L(e);return{key:s,first:null,c(){t=a("tr"),l=a("td"),c=a("div"),M.c(),f=p(),u=a("span"),m=b(_),q=p(),y=a("td"),g=a("span"),S=b(B),$=p(),R=a("td"),T&&T.c(),P=p(),w(c,"class","inline-flex"),w(g,"class","label"),this.first=t},m(k,F){r(k,t,F),i(t,l),i(l,c),M.m(c,null),i(c,f),i(c,u),i(u,m),i(t,q),i(t,y),i(y,g),i(g,S),i(t,$),i(t,R),T&&T.m(R,null),i(t,P)},p(k,F){e=k,D!==(D=I(e))&&(M.d(1),M=D(e),M&&(M.c(),M.m(c,f))),F&32&&_!==(_=e[15].name+"")&&ae(m,_),F&32&&B!==(B=we.getFieldValueType(e[15])+"")&&ae(S,B),L===(L=W(e))&&T?T.p(e,F):(T&&T.d(1),T=L&&L(e),T&&(T.c(),T.m(R,null)))},d(k){k&&o(t),M.d(),T&&T.d()}}}function wt(s,e){let t,l=e[10].code+"",c,f,u,_;function m(){return e[9](e[10])}return{key:s,first:null,c(){t=a("button"),c=b(l),f=p(),w(t,"class","tab-item"),Se(t,"active",e[2]===e[10].code),this.first=t},m(q,y){r(q,t,y),i(t,c),i(t,f),u||(_=Ht(t,"click",m),u=!0)},p(q,y){e=q,y&8&&l!==(l=e[10].code+"")&&ae(c,l),y&12&&Se(t,"active",e[2]===e[10].code)},d(q){q&&o(t),u=!1,_()}}}function gt(s,e){let t,l,c,f;return l=new Ct({props:{content:e[10].body}}),{key:s,first:null,c(){t=a("div"),Te(l.$$.fragment),c=p(),w(t,"class","tab-item"),Se(t,"active",e[2]===e[10].code),this.first=t},m(u,_){r(u,t,_),qe(l,t,null),i(t,c),f=!0},p(u,_){e=u;const m={};_&8&&(m.content=e[10].body),l.$set(m),(!f||_&12)&&Se(t,"active",e[2]===e[10].code)},i(u){f||(ve(l.$$.fragment,u),f=!0)},o(u){he(l.$$.fragment,u),f=!1},d(u){u&&o(t),$e(l)}}}function Kt(s){var st,at,ot,rt;let e,t,l=s[0].name+"",c,f,u,_,m,q,y,g=s[0].name+"",B,S,$,R,P,I,D,M,W,L,T,k,F,ee,z,U,oe,K,X=s[0].name+"",Y,ue,te,E,Z,G,H,h,C,N,x,O=[],je=new Map,Me,pe,Pe,le,Le,Je,me,ne,Fe,De,He,Ee,A,Ie,re,Ue,Qe,We,Oe,xe,Re,ze,Ke,Xe,Ae,Ye,Ze,de,Be,be,Ve,ie,_e,Q=[],Ge=new Map,et,ke,j=[],tt=new Map,se;M=new St({props:{js:`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${s[4]}');

...

// example create data
const data = ${JSON.stringify(s[7](s[0]),null,4)};

const record = await pb.collection('${(st=s[0])==null?void 0:st.name}').create(data);
`+(s[1]?`
// (optional) send an email verification request
await pb.collection('${(at=s[0])==null?void 0:at.name}').requestVerification('test@example.com');
`:""),dart:`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${s[4]}');

...

// example create body
final body = <String, dynamic>${JSON.stringify(s[7](s[0]),null,2)};

final record = await pb.collection('${(ot=s[0])==null?void 0:ot.name}').create(body: body);
`+(s[1]?`
// (optional) send an email verification request
await pb.collection('${(rt=s[0])==null?void 0:rt.name}').requestVerification('test@example.com');
`:"")}});let J=s[6]&&kt(),V=s[1]&&yt(s),ge=ce(s[5]);const lt=n=>n[15].name;for(let n=0;n<ge.length;n+=1){let d=_t(s,ge,n),v=lt(d);je.set(v,O[n]=vt(v,d))}re=new Ct({props:{content:"?expand=relField1,relField2.subRelField"}}),de=new Rt({});let Ce=ce(s[3]);const nt=n=>n[10].code;for(let n=0;n<Ce.length;n+=1){let d=bt(s,Ce,n),v=nt(d);Ge.set(v,Q[n]=wt(v,d))}let ye=ce(s[3]);const it=n=>n[10].code;for(let n=0;n<ye.length;n+=1){let d=mt(s,ye,n),v=it(d);tt.set(v,j[n]=gt(v,d))}return{c(){e=a("h3"),t=b("Create ("),c=b(l),f=b(")"),u=p(),_=a("div"),m=a("p"),q=b("Create a new "),y=a("strong"),B=b(g),S=b(" record."),$=p(),R=a("p"),R.innerHTML=`Body parameters could be sent as <code>application/json</code> or
        <code>multipart/form-data</code>.`,P=p(),I=a("p"),I.innerHTML=`File upload is supported only via <code>multipart/form-data</code>.
        <br/>
        For more info and examples you could check the detailed
        <a href="https://pocketbase.io/docs/files-handling" target="_blank" rel="noopener noreferrer">Files upload and handling docs
        </a>.`,D=p(),Te(M.$$.fragment),W=p(),L=a("h6"),L.textContent="API details",T=p(),k=a("div"),F=a("strong"),F.textContent="POST",ee=p(),z=a("div"),U=a("p"),oe=b("/api/collections/"),K=a("strong"),Y=b(X),ue=b("/records"),te=p(),J&&J.c(),E=p(),Z=a("div"),Z.textContent="Body Parameters",G=p(),H=a("table"),h=a("thead"),h.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr>',C=p(),N=a("tbody"),V&&V.c(),x=p();for(let n=0;n<O.length;n+=1)O[n].c();Me=p(),pe=a("div"),pe.textContent="Query parameters",Pe=p(),le=a("table"),Le=a("thead"),Le.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr>',Je=p(),me=a("tbody"),ne=a("tr"),Fe=a("td"),Fe.textContent="expand",De=p(),He=a("td"),He.innerHTML='<span class="label">String</span>',Ee=p(),A=a("td"),Ie=b(`Auto expand relations when returning the created record. Ex.:
                `),Te(re.$$.fragment),Ue=b(`
                Supports up to 6-levels depth nested relations expansion. `),Qe=a("br"),We=b(`
                The expanded relations will be appended to the record under the
                `),Oe=a("code"),Oe.textContent="expand",xe=b(" property (eg. "),Re=a("code"),Re.textContent='"expand": {"relField1": {...}, ...}',ze=b(`).
                `),Ke=a("br"),Xe=b(`
                Only the relations to which the request user has permissions to `),Ae=a("strong"),Ae.textContent="view",Ye=b(" will be expanded."),Ze=p(),Te(de.$$.fragment),Be=p(),be=a("div"),be.textContent="Responses",Ve=p(),ie=a("div"),_e=a("div");for(let n=0;n<Q.length;n+=1)Q[n].c();et=p(),ke=a("div");for(let n=0;n<j.length;n+=1)j[n].c();w(e,"class","m-b-sm"),w(_,"class","content txt-lg m-b-sm"),w(L,"class","m-b-xs"),w(F,"class","label label-primary"),w(z,"class","content"),w(k,"class","alert alert-success"),w(Z,"class","section-title"),w(H,"class","table-compact table-border m-b-base"),w(pe,"class","section-title"),w(le,"class","table-compact table-border m-b-base"),w(be,"class","section-title"),w(_e,"class","tabs-header compact combined left"),w(ke,"class","tabs-content"),w(ie,"class","tabs")},m(n,d){r(n,e,d),i(e,t),i(e,c),i(e,f),r(n,u,d),r(n,_,d),i(_,m),i(m,q),i(m,y),i(y,B),i(m,S),i(_,$),i(_,R),i(_,P),i(_,I),r(n,D,d),qe(M,n,d),r(n,W,d),r(n,L,d),r(n,T,d),r(n,k,d),i(k,F),i(k,ee),i(k,z),i(z,U),i(U,oe),i(U,K),i(K,Y),i(U,ue),i(k,te),J&&J.m(k,null),r(n,E,d),r(n,Z,d),r(n,G,d),r(n,H,d),i(H,h),i(H,C),i(H,N),V&&V.m(N,null),i(N,x);for(let v=0;v<O.length;v+=1)O[v]&&O[v].m(N,null);r(n,Me,d),r(n,pe,d),r(n,Pe,d),r(n,le,d),i(le,Le),i(le,Je),i(le,me),i(me,ne),i(ne,Fe),i(ne,De),i(ne,He),i(ne,Ee),i(ne,A),i(A,Ie),qe(re,A,null),i(A,Ue),i(A,Qe),i(A,We),i(A,Oe),i(A,xe),i(A,Re),i(A,ze),i(A,Ke),i(A,Xe),i(A,Ae),i(A,Ye),i(me,Ze),qe(de,me,null),r(n,Be,d),r(n,be,d),r(n,Ve,d),r(n,ie,d),i(ie,_e);for(let v=0;v<Q.length;v+=1)Q[v]&&Q[v].m(_e,null);i(ie,et),i(ie,ke);for(let v=0;v<j.length;v+=1)j[v]&&j[v].m(ke,null);se=!0},p(n,[d]){var dt,ct,ft,ut;(!se||d&1)&&l!==(l=n[0].name+"")&&ae(c,l),(!se||d&1)&&g!==(g=n[0].name+"")&&ae(B,g);const v={};d&19&&(v.js=`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${n[4]}');

...

// example create data
const data = ${JSON.stringify(n[7](n[0]),null,4)};

const record = await pb.collection('${(dt=n[0])==null?void 0:dt.name}').create(data);
`+(n[1]?`
// (optional) send an email verification request
await pb.collection('${(ct=n[0])==null?void 0:ct.name}').requestVerification('test@example.com');
`:"")),d&19&&(v.dart=`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${n[4]}');

...

// example create body
final body = <String, dynamic>${JSON.stringify(n[7](n[0]),null,2)};

final record = await pb.collection('${(ft=n[0])==null?void 0:ft.name}').create(body: body);
`+(n[1]?`
// (optional) send an email verification request
await pb.collection('${(ut=n[0])==null?void 0:ut.name}').requestVerification('test@example.com');
`:"")),M.$set(v),(!se||d&1)&&X!==(X=n[0].name+"")&&ae(Y,X),n[6]?J||(J=kt(),J.c(),J.m(k,null)):J&&(J.d(1),J=null),n[1]?V?V.p(n,d):(V=yt(n),V.c(),V.m(N,x)):V&&(V.d(1),V=null),d&32&&(ge=ce(n[5]),O=Ne(O,d,lt,1,n,ge,je,N,pt,vt,null,_t)),d&12&&(Ce=ce(n[3]),Q=Ne(Q,d,nt,1,n,Ce,Ge,_e,pt,wt,null,bt)),d&12&&(ye=ce(n[3]),Mt(),j=Ne(j,d,it,1,n,ye,tt,ke,Pt,gt,null,mt),Lt())},i(n){if(!se){ve(M.$$.fragment,n),ve(re.$$.fragment,n),ve(de.$$.fragment,n);for(let d=0;d<ye.length;d+=1)ve(j[d]);se=!0}},o(n){he(M.$$.fragment,n),he(re.$$.fragment,n),he(de.$$.fragment,n);for(let d=0;d<j.length;d+=1)he(j[d]);se=!1},d(n){n&&(o(e),o(u),o(_),o(D),o(W),o(L),o(T),o(k),o(E),o(Z),o(G),o(H),o(Me),o(pe),o(Pe),o(le),o(Be),o(be),o(Ve),o(ie)),$e(M,n),J&&J.d(),V&&V.d();for(let d=0;d<O.length;d+=1)O[d].d();$e(re),$e(de);for(let d=0;d<Q.length;d+=1)Q[d].d();for(let d=0;d<j.length;d+=1)j[d].d()}}}const Xt=s=>s.name=="emailVisibility",Yt=s=>s.name=="email";function Zt(s,e,t){let l,c,f,u,_,{collection:m}=e,q=200,y=[];function g(S){let $=we.dummyCollectionSchemaData(S,!0);return l&&($.password="12345678",$.passwordConfirm="12345678",delete $.verified),$}const B=S=>t(2,q=S.code);return s.$$set=S=>{"collection"in S&&t(0,m=S.collection)},s.$$.update=()=>{var S,$,R;s.$$.dirty&1&&t(1,l=m.type==="auth"),s.$$.dirty&1&&t(6,c=(m==null?void 0:m.createRule)===null),s.$$.dirty&2&&t(8,f=l?["password","verified","email","emailVisibility"]:[]),s.$$.dirty&257&&t(5,u=((S=m==null?void 0:m.fields)==null?void 0:S.filter(P=>!P.hidden&&P.type!="autodate"&&!f.includes(P.name)))||[]),s.$$.dirty&1&&t(3,y=[{code:200,body:JSON.stringify(we.dummyCollectionRecord(m),null,2)},{code:400,body:`
                {
                  "status": 400,
                  "message": "Failed to create record.",
                  "data": {
                    "${(R=($=m==null?void 0:m.fields)==null?void 0:$[0])==null?void 0:R.name}": {
                      "code": "validation_required",
                      "message": "Missing required value."
                    }
                  }
                }
            `},{code:403,body:`
                {
                  "status": 403,
                  "message": "You are not allowed to perform this request.",
                  "data": {}
                }
            `}])},t(4,_=we.getApiExampleUrl(Ft.baseURL)),[m,l,q,y,_,u,c,g,f,B]}class tl extends $t{constructor(e){super(),qt(this,e,Zt,Kt,Tt,{collection:0})}}export{tl as default};
