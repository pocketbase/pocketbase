import{S as $t,i as qt,s as Tt,V as St,X as ce,W as Ct,h as r,d as $e,t as he,a as ye,I as ae,Z as Ne,_ as pt,C as Mt,$ as Pt,D as Lt,l as d,n as i,m as qe,u as a,A as _,v as p,c as Te,w,J as ve,p as Ft,k as Se,o as Ht,L as Ot,H as we}from"./index-C8ya_RsB.js";import{F as Rt}from"./FieldsQueryParam-SLyti2ZR.js";function mt(s,e,t){const l=s.slice();return l[10]=e[t],l}function bt(s,e,t){const l=s.slice();return l[10]=e[t],l}function _t(s,e,t){const l=s.slice();return l[15]=e[t],l}function kt(s){let e;return{c(){e=a("p"),e.innerHTML="Requires superuser <code>Authorization:TOKEN</code> header",w(e,"class","txt-hint txt-sm txt-right")},m(t,l){d(t,e,l)},d(t){t&&r(e)}}}function ht(s){let e,t,l,f,c,u,b,m,q,h,g,B,S,$,R,P,I,D,M,W,L,T,k,F,ee,K,U,oe,X,Y,Z;function fe(y,C){var N,z,O;return C&1&&(u=null),u==null&&(u=!!((O=(z=(N=y[0])==null?void 0:N.fields)==null?void 0:z.find(Yt))!=null&&O.required)),u?Bt:At}let te=fe(s,-1),E=te(s);function G(y,C){var N,z,O;return C&1&&(I=null),I==null&&(I=!!((O=(z=(N=y[0])==null?void 0:N.fields)==null?void 0:z.find(Xt))!=null&&O.required)),I?Nt:Vt}let x=G(s,-1),H=x(s);return{c(){e=a("tr"),e.innerHTML='<td colspan="3" class="txt-hint txt-bold">Auth specific fields</td>',t=p(),l=a("tr"),f=a("td"),c=a("div"),E.c(),b=p(),m=a("span"),m.textContent="email",q=p(),h=a("td"),h.innerHTML='<span class="label">String</span>',g=p(),B=a("td"),B.textContent="Auth record email address.",S=p(),$=a("tr"),R=a("td"),P=a("div"),H.c(),D=p(),M=a("span"),M.textContent="emailVisibility",W=p(),L=a("td"),L.innerHTML='<span class="label">Boolean</span>',T=p(),k=a("td"),k.textContent="Whether to show/hide the auth record email when fetching the record data.",F=p(),ee=a("tr"),ee.innerHTML='<td><div class="inline-flex"><span class="label label-success">Required</span> <span>password</span></div></td> <td><span class="label">String</span></td> <td>Auth record password.</td>',K=p(),U=a("tr"),U.innerHTML='<td><div class="inline-flex"><span class="label label-success">Required</span> <span>passwordConfirm</span></div></td> <td><span class="label">String</span></td> <td>Auth record password confirmation.</td>',oe=p(),X=a("tr"),X.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>verified</span></div></td> <td><span class="label">Boolean</span></td> <td>Indicates whether the auth record is verified or not.
                    <br/>
                    This field can be set only by superusers or auth records with &quot;Manage&quot; access.</td>`,Y=p(),Z=a("tr"),Z.innerHTML='<td colspan="3" class="txt-hint txt-bold">Other fields</td>',w(c,"class","inline-flex"),w(P,"class","inline-flex")},m(y,C){d(y,e,C),d(y,t,C),d(y,l,C),i(l,f),i(f,c),E.m(c,null),i(c,b),i(c,m),i(l,q),i(l,h),i(l,g),i(l,B),d(y,S,C),d(y,$,C),i($,R),i(R,P),H.m(P,null),i(P,D),i(P,M),i($,W),i($,L),i($,T),i($,k),d(y,F,C),d(y,ee,C),d(y,K,C),d(y,U,C),d(y,oe,C),d(y,X,C),d(y,Y,C),d(y,Z,C)},p(y,C){te!==(te=fe(y,C))&&(E.d(1),E=te(y),E&&(E.c(),E.m(c,b))),x!==(x=G(y,C))&&(H.d(1),H=x(y),H&&(H.c(),H.m(P,D)))},d(y){y&&(r(e),r(t),r(l),r(S),r($),r(F),r(ee),r(K),r(U),r(oe),r(X),r(Y),r(Z)),E.d(),H.d()}}}function At(s){let e;return{c(){e=a("span"),e.textContent="Optional",w(e,"class","label label-warning")},m(t,l){d(t,e,l)},d(t){t&&r(e)}}}function Bt(s){let e;return{c(){e=a("span"),e.textContent="Required",w(e,"class","label label-success")},m(t,l){d(t,e,l)},d(t){t&&r(e)}}}function Vt(s){let e;return{c(){e=a("span"),e.textContent="Optional",w(e,"class","label label-warning")},m(t,l){d(t,e,l)},d(t){t&&r(e)}}}function Nt(s){let e;return{c(){e=a("span"),e.textContent="Required",w(e,"class","label label-success")},m(t,l){d(t,e,l)},d(t){t&&r(e)}}}function Jt(s){let e;return{c(){e=a("span"),e.textContent="Required",w(e,"class","label label-success")},m(t,l){d(t,e,l)},d(t){t&&r(e)}}}function jt(s){let e;return{c(){e=a("span"),e.textContent="Optional",w(e,"class","label label-warning")},m(t,l){d(t,e,l)},d(t){t&&r(e)}}}function Dt(s){let e,t=s[15].maxSelect===1?"id":"ids",l,f;return{c(){e=_("Relation record "),l=_(t),f=_(".")},m(c,u){d(c,e,u),d(c,l,u),d(c,f,u)},p(c,u){u&32&&t!==(t=c[15].maxSelect===1?"id":"ids")&&ae(l,t)},d(c){c&&(r(e),r(l),r(f))}}}function Et(s){let e,t,l,f,c,u,b,m,q;return{c(){e=_("File object."),t=a("br"),l=_(`
                        Set to empty value (`),f=a("code"),f.textContent="null",c=_(", "),u=a("code"),u.textContent='""',b=_(" or "),m=a("code"),m.textContent="[]",q=_(`) to delete
                        already uploaded file(s).`)},m(h,g){d(h,e,g),d(h,t,g),d(h,l,g),d(h,f,g),d(h,c,g),d(h,u,g),d(h,b,g),d(h,m,g),d(h,q,g)},p:we,d(h){h&&(r(e),r(t),r(l),r(f),r(c),r(u),r(b),r(m),r(q))}}}function It(s){let e;return{c(){e=_("URL address.")},m(t,l){d(t,e,l)},p:we,d(t){t&&r(e)}}}function Ut(s){let e;return{c(){e=_("Email address.")},m(t,l){d(t,e,l)},p:we,d(t){t&&r(e)}}}function Qt(s){let e;return{c(){e=_("JSON array or object.")},m(t,l){d(t,e,l)},p:we,d(t){t&&r(e)}}}function Wt(s){let e;return{c(){e=_("Number value.")},m(t,l){d(t,e,l)},p:we,d(t){t&&r(e)}}}function zt(s){let e,t,l=s[15].autogeneratePattern&&yt();return{c(){e=_(`Plain text value.
                        `),l&&l.c(),t=Ot()},m(f,c){d(f,e,c),l&&l.m(f,c),d(f,t,c)},p(f,c){f[15].autogeneratePattern?l||(l=yt(),l.c(),l.m(t.parentNode,t)):l&&(l.d(1),l=null)},d(f){f&&(r(e),r(t)),l&&l.d(f)}}}function yt(s){let e;return{c(){e=_("It is autogenerated if not set.")},m(t,l){d(t,e,l)},d(t){t&&r(e)}}}function vt(s,e){let t,l,f,c,u,b=e[15].name+"",m,q,h,g,B=ve.getFieldValueType(e[15])+"",S,$,R,P;function I(k,F){return!k[15].required||k[15].type=="text"&&k[15].autogeneratePattern?jt:Jt}let D=I(e),M=D(e);function W(k,F){if(k[15].type==="text")return zt;if(k[15].type==="number")return Wt;if(k[15].type==="json")return Qt;if(k[15].type==="email")return Ut;if(k[15].type==="url")return It;if(k[15].type==="file")return Et;if(k[15].type==="relation")return Dt}let L=W(e),T=L&&L(e);return{key:s,first:null,c(){t=a("tr"),l=a("td"),f=a("div"),M.c(),c=p(),u=a("span"),m=_(b),q=p(),h=a("td"),g=a("span"),S=_(B),$=p(),R=a("td"),T&&T.c(),P=p(),w(f,"class","inline-flex"),w(g,"class","label"),this.first=t},m(k,F){d(k,t,F),i(t,l),i(l,f),M.m(f,null),i(f,c),i(f,u),i(u,m),i(t,q),i(t,h),i(h,g),i(g,S),i(t,$),i(t,R),T&&T.m(R,null),i(t,P)},p(k,F){e=k,D!==(D=I(e))&&(M.d(1),M=D(e),M&&(M.c(),M.m(f,c))),F&32&&b!==(b=e[15].name+"")&&ae(m,b),F&32&&B!==(B=ve.getFieldValueType(e[15])+"")&&ae(S,B),L===(L=W(e))&&T?T.p(e,F):(T&&T.d(1),T=L&&L(e),T&&(T.c(),T.m(R,null)))},d(k){k&&r(t),M.d(),T&&T.d()}}}function wt(s,e){let t,l=e[10].code+"",f,c,u,b;function m(){return e[9](e[10])}return{key:s,first:null,c(){t=a("button"),f=_(l),c=p(),w(t,"class","tab-item"),Se(t,"active",e[2]===e[10].code),this.first=t},m(q,h){d(q,t,h),i(t,f),i(t,c),u||(b=Ht(t,"click",m),u=!0)},p(q,h){e=q,h&8&&l!==(l=e[10].code+"")&&ae(f,l),h&12&&Se(t,"active",e[2]===e[10].code)},d(q){q&&r(t),u=!1,b()}}}function gt(s,e){let t,l,f,c;return l=new Ct({props:{content:e[10].body}}),{key:s,first:null,c(){t=a("div"),Te(l.$$.fragment),f=p(),w(t,"class","tab-item"),Se(t,"active",e[2]===e[10].code),this.first=t},m(u,b){d(u,t,b),qe(l,t,null),i(t,f),c=!0},p(u,b){e=u;const m={};b&8&&(m.content=e[10].body),l.$set(m),(!c||b&12)&&Se(t,"active",e[2]===e[10].code)},i(u){c||(ye(l.$$.fragment,u),c=!0)},o(u){he(l.$$.fragment,u),c=!1},d(u){u&&r(t),$e(l)}}}function Kt(s){var st,at,ot,rt;let e,t,l=s[0].name+"",f,c,u,b,m,q,h,g=s[0].name+"",B,S,$,R,P,I,D,M,W,L,T,k,F,ee,K,U,oe,X,Y=s[0].name+"",Z,fe,te,E,G,x,H,y,C,N,z,O=[],Je=new Map,Me,ue,Pe,le,Le,je,pe,ne,Fe,De,He,Ee,A,Ie,re,Ue,Qe,We,Oe,ze,Re,Ke,Xe,Ye,Ae,Ze,Ge,de,Be,me,Ve,ie,be,Q=[],xe=new Map,et,_e,J=[],tt=new Map,se;M=new St({props:{js:`
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
`:"")}});let j=s[6]&&kt(),V=s[1]&&ht(s),ge=ce(s[5]);const lt=n=>n[15].name;for(let n=0;n<ge.length;n+=1){let o=_t(s,ge,n),v=lt(o);Je.set(v,O[n]=vt(v,o))}re=new Ct({props:{content:"?expand=relField1,relField2.subRelField"}}),de=new Rt({});let Ce=ce(s[3]);const nt=n=>n[10].code;for(let n=0;n<Ce.length;n+=1){let o=bt(s,Ce,n),v=nt(o);xe.set(v,Q[n]=wt(v,o))}let ke=ce(s[3]);const it=n=>n[10].code;for(let n=0;n<ke.length;n+=1){let o=mt(s,ke,n),v=it(o);tt.set(v,J[n]=gt(v,o))}return{c(){e=a("h3"),t=_("Create ("),f=_(l),c=_(")"),u=p(),b=a("div"),m=a("p"),q=_("Create a new "),h=a("strong"),B=_(g),S=_(" record."),$=p(),R=a("p"),R.innerHTML=`Body parameters could be sent as <code>application/json</code> or
        <code>multipart/form-data</code>.`,P=p(),I=a("p"),I.innerHTML=`File upload is supported only via <code>multipart/form-data</code>.
        <br/>
        For more info and examples you could check the detailed
        <a href="https://pocketbase.io/docs/files-handling" target="_blank" rel="noopener noreferrer">Files upload and handling docs
        </a>.`,D=p(),Te(M.$$.fragment),W=p(),L=a("h6"),L.textContent="API details",T=p(),k=a("div"),F=a("strong"),F.textContent="POST",ee=p(),K=a("div"),U=a("p"),oe=_("/api/collections/"),X=a("strong"),Z=_(Y),fe=_("/records"),te=p(),j&&j.c(),E=p(),G=a("div"),G.textContent="Body Parameters",x=p(),H=a("table"),y=a("thead"),y.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr>',C=p(),N=a("tbody"),V&&V.c(),z=p();for(let n=0;n<O.length;n+=1)O[n].c();Me=p(),ue=a("div"),ue.textContent="Query parameters",Pe=p(),le=a("table"),Le=a("thead"),Le.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr>',je=p(),pe=a("tbody"),ne=a("tr"),Fe=a("td"),Fe.textContent="expand",De=p(),He=a("td"),He.innerHTML='<span class="label">String</span>',Ee=p(),A=a("td"),Ie=_(`Auto expand relations when returning the created record. Ex.:
                `),Te(re.$$.fragment),Ue=_(`
                Supports up to 6-levels depth nested relations expansion. `),Qe=a("br"),We=_(`
                The expanded relations will be appended to the record under the
                `),Oe=a("code"),Oe.textContent="expand",ze=_(" property (eg. "),Re=a("code"),Re.textContent='"expand": {"relField1": {...}, ...}',Ke=_(`).
                `),Xe=a("br"),Ye=_(`
                Only the relations to which the request user has permissions to `),Ae=a("strong"),Ae.textContent="view",Ze=_(" will be expanded."),Ge=p(),Te(de.$$.fragment),Be=p(),me=a("div"),me.textContent="Responses",Ve=p(),ie=a("div"),be=a("div");for(let n=0;n<Q.length;n+=1)Q[n].c();et=p(),_e=a("div");for(let n=0;n<J.length;n+=1)J[n].c();w(e,"class","m-b-sm"),w(b,"class","content txt-lg m-b-sm"),w(L,"class","m-b-xs"),w(F,"class","label label-primary"),w(K,"class","content"),w(k,"class","alert alert-success"),w(G,"class","section-title"),w(H,"class","table-compact table-border m-b-base"),w(ue,"class","section-title"),w(le,"class","table-compact table-border m-b-base"),w(me,"class","section-title"),w(be,"class","tabs-header compact combined left"),w(_e,"class","tabs-content"),w(ie,"class","tabs")},m(n,o){d(n,e,o),i(e,t),i(e,f),i(e,c),d(n,u,o),d(n,b,o),i(b,m),i(m,q),i(m,h),i(h,B),i(m,S),i(b,$),i(b,R),i(b,P),i(b,I),d(n,D,o),qe(M,n,o),d(n,W,o),d(n,L,o),d(n,T,o),d(n,k,o),i(k,F),i(k,ee),i(k,K),i(K,U),i(U,oe),i(U,X),i(X,Z),i(U,fe),i(k,te),j&&j.m(k,null),d(n,E,o),d(n,G,o),d(n,x,o),d(n,H,o),i(H,y),i(H,C),i(H,N),V&&V.m(N,null),i(N,z);for(let v=0;v<O.length;v+=1)O[v]&&O[v].m(N,null);d(n,Me,o),d(n,ue,o),d(n,Pe,o),d(n,le,o),i(le,Le),i(le,je),i(le,pe),i(pe,ne),i(ne,Fe),i(ne,De),i(ne,He),i(ne,Ee),i(ne,A),i(A,Ie),qe(re,A,null),i(A,Ue),i(A,Qe),i(A,We),i(A,Oe),i(A,ze),i(A,Re),i(A,Ke),i(A,Xe),i(A,Ye),i(A,Ae),i(A,Ze),i(pe,Ge),qe(de,pe,null),d(n,Be,o),d(n,me,o),d(n,Ve,o),d(n,ie,o),i(ie,be);for(let v=0;v<Q.length;v+=1)Q[v]&&Q[v].m(be,null);i(ie,et),i(ie,_e);for(let v=0;v<J.length;v+=1)J[v]&&J[v].m(_e,null);se=!0},p(n,[o]){var dt,ct,ft,ut;(!se||o&1)&&l!==(l=n[0].name+"")&&ae(f,l),(!se||o&1)&&g!==(g=n[0].name+"")&&ae(B,g);const v={};o&19&&(v.js=`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${n[4]}');

...

// example create data
const data = ${JSON.stringify(n[7](n[0]),null,4)};

const record = await pb.collection('${(dt=n[0])==null?void 0:dt.name}').create(data);
`+(n[1]?`
// (optional) send an email verification request
await pb.collection('${(ct=n[0])==null?void 0:ct.name}').requestVerification('test@example.com');
`:"")),o&19&&(v.dart=`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${n[4]}');

...

// example create body
final body = <String, dynamic>${JSON.stringify(n[7](n[0]),null,2)};

final record = await pb.collection('${(ft=n[0])==null?void 0:ft.name}').create(body: body);
`+(n[1]?`
// (optional) send an email verification request
await pb.collection('${(ut=n[0])==null?void 0:ut.name}').requestVerification('test@example.com');
`:"")),M.$set(v),(!se||o&1)&&Y!==(Y=n[0].name+"")&&ae(Z,Y),n[6]?j||(j=kt(),j.c(),j.m(k,null)):j&&(j.d(1),j=null),n[1]?V?V.p(n,o):(V=ht(n),V.c(),V.m(N,z)):V&&(V.d(1),V=null),o&32&&(ge=ce(n[5]),O=Ne(O,o,lt,1,n,ge,Je,N,pt,vt,null,_t)),o&12&&(Ce=ce(n[3]),Q=Ne(Q,o,nt,1,n,Ce,xe,be,pt,wt,null,bt)),o&12&&(ke=ce(n[3]),Mt(),J=Ne(J,o,it,1,n,ke,tt,_e,Pt,gt,null,mt),Lt())},i(n){if(!se){ye(M.$$.fragment,n),ye(re.$$.fragment,n),ye(de.$$.fragment,n);for(let o=0;o<ke.length;o+=1)ye(J[o]);se=!0}},o(n){he(M.$$.fragment,n),he(re.$$.fragment,n),he(de.$$.fragment,n);for(let o=0;o<J.length;o+=1)he(J[o]);se=!1},d(n){n&&(r(e),r(u),r(b),r(D),r(W),r(L),r(T),r(k),r(E),r(G),r(x),r(H),r(Me),r(ue),r(Pe),r(le),r(Be),r(me),r(Ve),r(ie)),$e(M,n),j&&j.d(),V&&V.d();for(let o=0;o<O.length;o+=1)O[o].d();$e(re),$e(de);for(let o=0;o<Q.length;o+=1)Q[o].d();for(let o=0;o<J.length;o+=1)J[o].d()}}}const Xt=s=>s.name=="emailVisibility",Yt=s=>s.name=="email";function Zt(s,e,t){let l,f,c,u,b,{collection:m}=e,q=200,h=[];function g(S){let $=ve.dummyCollectionSchemaData(S,!0);return l&&($.password="12345678",$.passwordConfirm="12345678",delete $.verified),$}const B=S=>t(2,q=S.code);return s.$$set=S=>{"collection"in S&&t(0,m=S.collection)},s.$$.update=()=>{var S,$,R;s.$$.dirty&1&&t(1,l=m.type==="auth"),s.$$.dirty&1&&t(6,f=(m==null?void 0:m.createRule)===null),s.$$.dirty&2&&t(8,c=l?["password","verified","email","emailVisibility"]:[]),s.$$.dirty&257&&t(5,u=((S=m==null?void 0:m.fields)==null?void 0:S.filter(P=>!P.hidden&&P.type!="autodate"&&!c.includes(P.name)))||[]),s.$$.dirty&1&&t(3,h=[{code:200,body:JSON.stringify(ve.dummyCollectionRecord(m),null,2)},{code:400,body:`
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
            `}])},t(4,b=ve.getApiExampleUrl(Ft.baseURL)),[m,l,q,h,b,u,f,g,c,B]}class el extends $t{constructor(e){super(),qt(this,e,Zt,Kt,Tt,{collection:0})}}export{el as default};
