import{S as $t,i as qt,s as Tt,ab as St,O as ce,ac as Ct,h as o,d as $e,t as he,a as ve,I as se,ad as Ve,ae as pt,C as Mt,af as Pt,D as Ot,l as r,n as i,m as qe,u as s,A as m,v as p,c as Te,w,J as we,p as Ft,k as Se,o as Ht,N as Lt,H as fe}from"./index-B4ZsHsKR.js";import{F as Rt}from"./FieldsQueryParam-K1y4zYh0.js";function bt(a,e,t){const l=a.slice();return l[10]=e[t],l}function mt(a,e,t){const l=a.slice();return l[10]=e[t],l}function _t(a,e,t){const l=a.slice();return l[15]=e[t],l}function kt(a){let e;return{c(){e=s("p"),e.innerHTML="Requires superuser <code>Authorization:TOKEN</code> header",w(e,"class","txt-hint txt-sm txt-right")},m(t,l){r(t,e,l)},d(t){t&&o(e)}}}function yt(a){let e,t,l,c,f,u,_,b,q,y,g,B,S,$,R,P,I,D,M,x,O,T,k,F,ee,K,U,oe,W,Y,G;function ue(h,C){var V,z,L;return C&1&&(u=null),u==null&&(u=!!((L=(z=(V=h[0])==null?void 0:V.fields)==null?void 0:z.find(Gt))!=null&&L.required)),u?Bt:At}let te=ue(a,-1),E=te(a);function X(h,C){var V,z,L;return C&1&&(I=null),I==null&&(I=!!((L=(z=(V=h[0])==null?void 0:V.fields)==null?void 0:z.find(Yt))!=null&&L.required)),I?Vt:Nt}let Z=X(a,-1),H=Z(a);return{c(){e=s("tr"),e.innerHTML='<td colspan="3" class="txt-hint txt-bold">Auth specific fields</td>',t=p(),l=s("tr"),c=s("td"),f=s("div"),E.c(),_=p(),b=s("span"),b.textContent="email",q=p(),y=s("td"),y.innerHTML='<span class="label">String</span>',g=p(),B=s("td"),B.textContent="Auth record email address.",S=p(),$=s("tr"),R=s("td"),P=s("div"),H.c(),D=p(),M=s("span"),M.textContent="emailVisibility",x=p(),O=s("td"),O.innerHTML='<span class="label">Boolean</span>',T=p(),k=s("td"),k.textContent="Whether to show/hide the auth record email when fetching the record data.",F=p(),ee=s("tr"),ee.innerHTML='<td><div class="inline-flex"><span class="label label-success">Required</span> <span>password</span></div></td> <td><span class="label">String</span></td> <td>Auth record password.</td>',K=p(),U=s("tr"),U.innerHTML='<td><div class="inline-flex"><span class="label label-success">Required</span> <span>passwordConfirm</span></div></td> <td><span class="label">String</span></td> <td>Auth record password confirmation.</td>',oe=p(),W=s("tr"),W.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>verified</span></div></td> <td><span class="label">Boolean</span></td> <td>Indicates whether the auth record is verified or not.
                    <br/>
                    This field can be set only by superusers or auth records with &quot;Manage&quot; access.</td>`,Y=p(),G=s("tr"),G.innerHTML='<td colspan="3" class="txt-hint txt-bold">Other fields</td>',w(f,"class","inline-flex"),w(P,"class","inline-flex")},m(h,C){r(h,e,C),r(h,t,C),r(h,l,C),i(l,c),i(c,f),E.m(f,null),i(f,_),i(f,b),i(l,q),i(l,y),i(l,g),i(l,B),r(h,S,C),r(h,$,C),i($,R),i(R,P),H.m(P,null),i(P,D),i(P,M),i($,x),i($,O),i($,T),i($,k),r(h,F,C),r(h,ee,C),r(h,K,C),r(h,U,C),r(h,oe,C),r(h,W,C),r(h,Y,C),r(h,G,C)},p(h,C){te!==(te=ue(h,C))&&(E.d(1),E=te(h),E&&(E.c(),E.m(f,_))),Z!==(Z=X(h,C))&&(H.d(1),H=Z(h),H&&(H.c(),H.m(P,D)))},d(h){h&&(o(e),o(t),o(l),o(S),o($),o(F),o(ee),o(K),o(U),o(oe),o(W),o(Y),o(G)),E.d(),H.d()}}}function At(a){let e;return{c(){e=s("span"),e.textContent="Optional",w(e,"class","label label-warning")},m(t,l){r(t,e,l)},d(t){t&&o(e)}}}function Bt(a){let e;return{c(){e=s("span"),e.textContent="Required",w(e,"class","label label-success")},m(t,l){r(t,e,l)},d(t){t&&o(e)}}}function Nt(a){let e;return{c(){e=s("span"),e.textContent="Optional",w(e,"class","label label-warning")},m(t,l){r(t,e,l)},d(t){t&&o(e)}}}function Vt(a){let e;return{c(){e=s("span"),e.textContent="Required",w(e,"class","label label-success")},m(t,l){r(t,e,l)},d(t){t&&o(e)}}}function jt(a){let e;return{c(){e=s("span"),e.textContent="Required",w(e,"class","label label-success")},m(t,l){r(t,e,l)},d(t){t&&o(e)}}}function Jt(a){let e;return{c(){e=s("span"),e.textContent="Optional",w(e,"class","label label-warning")},m(t,l){r(t,e,l)},d(t){t&&o(e)}}}function Dt(a){let e,t=a[15].maxSelect===1?"id":"ids",l,c;return{c(){e=m("Relation record "),l=m(t),c=m(".")},m(f,u){r(f,e,u),r(f,l,u),r(f,c,u)},p(f,u){u&32&&t!==(t=f[15].maxSelect===1?"id":"ids")&&se(l,t)},d(f){f&&(o(e),o(l),o(c))}}}function Et(a){let e,t,l,c,f,u,_,b,q;return{c(){e=m("File object."),t=s("br"),l=m(`
                        Set to empty value (`),c=s("code"),c.textContent="null",f=m(", "),u=s("code"),u.textContent='""',_=m(" or "),b=s("code"),b.textContent="[]",q=m(`) to delete
                        already uploaded file(s).`)},m(y,g){r(y,e,g),r(y,t,g),r(y,l,g),r(y,c,g),r(y,f,g),r(y,u,g),r(y,_,g),r(y,b,g),r(y,q,g)},p:fe,d(y){y&&(o(e),o(t),o(l),o(c),o(f),o(u),o(_),o(b),o(q))}}}function It(a){let e,t;return{c(){e=s("code"),e.textContent='{"lon":x,"lat":y}',t=m(" object.")},m(l,c){r(l,e,c),r(l,t,c)},p:fe,d(l){l&&(o(e),o(t))}}}function Ut(a){let e;return{c(){e=m("URL address.")},m(t,l){r(t,e,l)},p:fe,d(t){t&&o(e)}}}function Qt(a){let e;return{c(){e=m("Email address.")},m(t,l){r(t,e,l)},p:fe,d(t){t&&o(e)}}}function xt(a){let e;return{c(){e=m("JSON array or object.")},m(t,l){r(t,e,l)},p:fe,d(t){t&&o(e)}}}function zt(a){let e;return{c(){e=m("Number value.")},m(t,l){r(t,e,l)},p:fe,d(t){t&&o(e)}}}function Kt(a){let e,t,l=a[15].autogeneratePattern&&ht();return{c(){e=m(`Plain text value.
                        `),l&&l.c(),t=Lt()},m(c,f){r(c,e,f),l&&l.m(c,f),r(c,t,f)},p(c,f){c[15].autogeneratePattern?l||(l=ht(),l.c(),l.m(t.parentNode,t)):l&&(l.d(1),l=null)},d(c){c&&(o(e),o(t)),l&&l.d(c)}}}function ht(a){let e;return{c(){e=m("It is autogenerated if not set.")},m(t,l){r(t,e,l)},d(t){t&&o(e)}}}function vt(a,e){let t,l,c,f,u,_=e[15].name+"",b,q,y,g,B=we.getFieldValueType(e[15])+"",S,$,R,P;function I(k,F){return!k[15].required||k[15].type=="text"&&k[15].autogeneratePattern?Jt:jt}let D=I(e),M=D(e);function x(k,F){if(k[15].type==="text")return Kt;if(k[15].type==="number")return zt;if(k[15].type==="json")return xt;if(k[15].type==="email")return Qt;if(k[15].type==="url")return Ut;if(k[15].type==="geoPoint")return It;if(k[15].type==="file")return Et;if(k[15].type==="relation")return Dt}let O=x(e),T=O&&O(e);return{key:a,first:null,c(){t=s("tr"),l=s("td"),c=s("div"),M.c(),f=p(),u=s("span"),b=m(_),q=p(),y=s("td"),g=s("span"),S=m(B),$=p(),R=s("td"),T&&T.c(),P=p(),w(c,"class","inline-flex"),w(g,"class","label"),this.first=t},m(k,F){r(k,t,F),i(t,l),i(l,c),M.m(c,null),i(c,f),i(c,u),i(u,b),i(t,q),i(t,y),i(y,g),i(g,S),i(t,$),i(t,R),T&&T.m(R,null),i(t,P)},p(k,F){e=k,D!==(D=I(e))&&(M.d(1),M=D(e),M&&(M.c(),M.m(c,f))),F&32&&_!==(_=e[15].name+"")&&se(b,_),F&32&&B!==(B=we.getFieldValueType(e[15])+"")&&se(S,B),O===(O=x(e))&&T?T.p(e,F):(T&&T.d(1),T=O&&O(e),T&&(T.c(),T.m(R,null)))},d(k){k&&o(t),M.d(),T&&T.d()}}}function wt(a,e){let t,l=e[10].code+"",c,f,u,_;function b(){return e[9](e[10])}return{key:a,first:null,c(){t=s("button"),c=m(l),f=p(),w(t,"class","tab-item"),Se(t,"active",e[2]===e[10].code),this.first=t},m(q,y){r(q,t,y),i(t,c),i(t,f),u||(_=Ht(t,"click",b),u=!0)},p(q,y){e=q,y&8&&l!==(l=e[10].code+"")&&se(c,l),y&12&&Se(t,"active",e[2]===e[10].code)},d(q){q&&o(t),u=!1,_()}}}function gt(a,e){let t,l,c,f;return l=new Ct({props:{content:e[10].body}}),{key:a,first:null,c(){t=s("div"),Te(l.$$.fragment),c=p(),w(t,"class","tab-item"),Se(t,"active",e[2]===e[10].code),this.first=t},m(u,_){r(u,t,_),qe(l,t,null),i(t,c),f=!0},p(u,_){e=u;const b={};_&8&&(b.content=e[10].body),l.$set(b),(!f||_&12)&&Se(t,"active",e[2]===e[10].code)},i(u){f||(ve(l.$$.fragment,u),f=!0)},o(u){he(l.$$.fragment,u),f=!1},d(u){u&&o(t),$e(l)}}}function Wt(a){var at,st,ot,rt;let e,t,l=a[0].name+"",c,f,u,_,b,q,y,g=a[0].name+"",B,S,$,R,P,I,D,M,x,O,T,k,F,ee,K,U,oe,W,Y=a[0].name+"",G,ue,te,E,X,Z,H,h,C,V,z,L=[],je=new Map,Me,pe,Pe,le,Oe,Je,be,ne,Fe,De,He,Ee,A,Ie,re,Ue,Qe,xe,Le,ze,Re,Ke,We,Ye,Ae,Ge,Xe,de,Be,me,Ne,ie,_e,Q=[],Ze=new Map,et,ke,j=[],tt=new Map,ae;M=new St({props:{js:`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${a[4]}');

...

// example create data
const data = ${JSON.stringify(a[7](a[0]),null,4)};

const record = await pb.collection('${(at=a[0])==null?void 0:at.name}').create(data);
`+(a[1]?`
// (optional) send an email verification request
await pb.collection('${(st=a[0])==null?void 0:st.name}').requestVerification('test@example.com');
`:""),dart:`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${a[4]}');

...

// example create body
final body = <String, dynamic>${JSON.stringify(a[7](a[0]),null,2)};

final record = await pb.collection('${(ot=a[0])==null?void 0:ot.name}').create(body: body);
`+(a[1]?`
// (optional) send an email verification request
await pb.collection('${(rt=a[0])==null?void 0:rt.name}').requestVerification('test@example.com');
`:"")}});let J=a[6]&&kt(),N=a[1]&&yt(a),ge=ce(a[5]);const lt=n=>n[15].name;for(let n=0;n<ge.length;n+=1){let d=_t(a,ge,n),v=lt(d);je.set(v,L[n]=vt(v,d))}re=new Ct({props:{content:"?expand=relField1,relField2.subRelField"}}),de=new Rt({});let Ce=ce(a[3]);const nt=n=>n[10].code;for(let n=0;n<Ce.length;n+=1){let d=mt(a,Ce,n),v=nt(d);Ze.set(v,Q[n]=wt(v,d))}let ye=ce(a[3]);const it=n=>n[10].code;for(let n=0;n<ye.length;n+=1){let d=bt(a,ye,n),v=it(d);tt.set(v,j[n]=gt(v,d))}return{c(){e=s("h3"),t=m("Create ("),c=m(l),f=m(")"),u=p(),_=s("div"),b=s("p"),q=m("Create a new "),y=s("strong"),B=m(g),S=m(" record."),$=p(),R=s("p"),R.innerHTML=`Body parameters could be sent as <code>application/json</code> or
        <code>multipart/form-data</code>.`,P=p(),I=s("p"),I.innerHTML=`File upload is supported only via <code>multipart/form-data</code>.
        <br/>
        For more info and examples you could check the detailed
        <a href="https://pocketbase.io/docs/files-handling" target="_blank" rel="noopener noreferrer">Files upload and handling docs
        </a>.`,D=p(),Te(M.$$.fragment),x=p(),O=s("h6"),O.textContent="API details",T=p(),k=s("div"),F=s("strong"),F.textContent="POST",ee=p(),K=s("div"),U=s("p"),oe=m("/api/collections/"),W=s("strong"),G=m(Y),ue=m("/records"),te=p(),J&&J.c(),E=p(),X=s("div"),X.textContent="Body Parameters",Z=p(),H=s("table"),h=s("thead"),h.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr>',C=p(),V=s("tbody"),N&&N.c(),z=p();for(let n=0;n<L.length;n+=1)L[n].c();Me=p(),pe=s("div"),pe.textContent="Query parameters",Pe=p(),le=s("table"),Oe=s("thead"),Oe.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr>',Je=p(),be=s("tbody"),ne=s("tr"),Fe=s("td"),Fe.textContent="expand",De=p(),He=s("td"),He.innerHTML='<span class="label">String</span>',Ee=p(),A=s("td"),Ie=m(`Auto expand relations when returning the created record. Ex.:
                `),Te(re.$$.fragment),Ue=m(`
                Supports up to 6-levels depth nested relations expansion. `),Qe=s("br"),xe=m(`
                The expanded relations will be appended to the record under the
                `),Le=s("code"),Le.textContent="expand",ze=m(" property (eg. "),Re=s("code"),Re.textContent='"expand": {"relField1": {...}, ...}',Ke=m(`).
                `),We=s("br"),Ye=m(`
                Only the relations to which the request user has permissions to `),Ae=s("strong"),Ae.textContent="view",Ge=m(" will be expanded."),Xe=p(),Te(de.$$.fragment),Be=p(),me=s("div"),me.textContent="Responses",Ne=p(),ie=s("div"),_e=s("div");for(let n=0;n<Q.length;n+=1)Q[n].c();et=p(),ke=s("div");for(let n=0;n<j.length;n+=1)j[n].c();w(e,"class","m-b-sm"),w(_,"class","content txt-lg m-b-sm"),w(O,"class","m-b-xs"),w(F,"class","label label-primary"),w(K,"class","content"),w(k,"class","alert alert-success"),w(X,"class","section-title"),w(H,"class","table-compact table-border m-b-base"),w(pe,"class","section-title"),w(le,"class","table-compact table-border m-b-base"),w(me,"class","section-title"),w(_e,"class","tabs-header compact combined left"),w(ke,"class","tabs-content"),w(ie,"class","tabs")},m(n,d){r(n,e,d),i(e,t),i(e,c),i(e,f),r(n,u,d),r(n,_,d),i(_,b),i(b,q),i(b,y),i(y,B),i(b,S),i(_,$),i(_,R),i(_,P),i(_,I),r(n,D,d),qe(M,n,d),r(n,x,d),r(n,O,d),r(n,T,d),r(n,k,d),i(k,F),i(k,ee),i(k,K),i(K,U),i(U,oe),i(U,W),i(W,G),i(U,ue),i(k,te),J&&J.m(k,null),r(n,E,d),r(n,X,d),r(n,Z,d),r(n,H,d),i(H,h),i(H,C),i(H,V),N&&N.m(V,null),i(V,z);for(let v=0;v<L.length;v+=1)L[v]&&L[v].m(V,null);r(n,Me,d),r(n,pe,d),r(n,Pe,d),r(n,le,d),i(le,Oe),i(le,Je),i(le,be),i(be,ne),i(ne,Fe),i(ne,De),i(ne,He),i(ne,Ee),i(ne,A),i(A,Ie),qe(re,A,null),i(A,Ue),i(A,Qe),i(A,xe),i(A,Le),i(A,ze),i(A,Re),i(A,Ke),i(A,We),i(A,Ye),i(A,Ae),i(A,Ge),i(be,Xe),qe(de,be,null),r(n,Be,d),r(n,me,d),r(n,Ne,d),r(n,ie,d),i(ie,_e);for(let v=0;v<Q.length;v+=1)Q[v]&&Q[v].m(_e,null);i(ie,et),i(ie,ke);for(let v=0;v<j.length;v+=1)j[v]&&j[v].m(ke,null);ae=!0},p(n,[d]){var dt,ct,ft,ut;(!ae||d&1)&&l!==(l=n[0].name+"")&&se(c,l),(!ae||d&1)&&g!==(g=n[0].name+"")&&se(B,g);const v={};d&19&&(v.js=`
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
`:"")),M.$set(v),(!ae||d&1)&&Y!==(Y=n[0].name+"")&&se(G,Y),n[6]?J||(J=kt(),J.c(),J.m(k,null)):J&&(J.d(1),J=null),n[1]?N?N.p(n,d):(N=yt(n),N.c(),N.m(V,z)):N&&(N.d(1),N=null),d&32&&(ge=ce(n[5]),L=Ve(L,d,lt,1,n,ge,je,V,pt,vt,null,_t)),d&12&&(Ce=ce(n[3]),Q=Ve(Q,d,nt,1,n,Ce,Ze,_e,pt,wt,null,mt)),d&12&&(ye=ce(n[3]),Mt(),j=Ve(j,d,it,1,n,ye,tt,ke,Pt,gt,null,bt),Ot())},i(n){if(!ae){ve(M.$$.fragment,n),ve(re.$$.fragment,n),ve(de.$$.fragment,n);for(let d=0;d<ye.length;d+=1)ve(j[d]);ae=!0}},o(n){he(M.$$.fragment,n),he(re.$$.fragment,n),he(de.$$.fragment,n);for(let d=0;d<j.length;d+=1)he(j[d]);ae=!1},d(n){n&&(o(e),o(u),o(_),o(D),o(x),o(O),o(T),o(k),o(E),o(X),o(Z),o(H),o(Me),o(pe),o(Pe),o(le),o(Be),o(me),o(Ne),o(ie)),$e(M,n),J&&J.d(),N&&N.d();for(let d=0;d<L.length;d+=1)L[d].d();$e(re),$e(de);for(let d=0;d<Q.length;d+=1)Q[d].d();for(let d=0;d<j.length;d+=1)j[d].d()}}}const Yt=a=>a.name=="emailVisibility",Gt=a=>a.name=="email";function Xt(a,e,t){let l,c,f,u,_,{collection:b}=e,q=200,y=[];function g(S){let $=we.dummyCollectionSchemaData(S,!0);return l&&($.password="12345678",$.passwordConfirm="12345678",delete $.verified),$}const B=S=>t(2,q=S.code);return a.$$set=S=>{"collection"in S&&t(0,b=S.collection)},a.$$.update=()=>{var S,$,R;a.$$.dirty&1&&t(1,l=b.type==="auth"),a.$$.dirty&1&&t(6,c=(b==null?void 0:b.createRule)===null),a.$$.dirty&2&&t(8,f=l?["password","verified","email","emailVisibility"]:[]),a.$$.dirty&257&&t(5,u=((S=b==null?void 0:b.fields)==null?void 0:S.filter(P=>!P.hidden&&P.type!="autodate"&&!f.includes(P.name)))||[]),a.$$.dirty&1&&t(3,y=[{code:200,body:JSON.stringify(we.dummyCollectionRecord(b),null,2)},{code:400,body:`
                {
                  "status": 400,
                  "message": "Failed to create record.",
                  "data": {
                    "${(R=($=b==null?void 0:b.fields)==null?void 0:$[0])==null?void 0:R.name}": {
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
            `}])},t(4,_=we.getApiExampleUrl(Ft.baseURL)),[b,l,q,y,_,u,c,g,f,B]}class tl extends $t{constructor(e){super(),qt(this,e,Xt,Wt,Tt,{collection:0})}}export{tl as default};
