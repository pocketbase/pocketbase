import{S as $t,i as qt,s as St,V as Tt,J as ee,X as ue,W as Ct,h as s,z as _,j as p,c as $e,k as w,n as r,o as i,m as qe,H as oe,Y as Ve,Z as pt,E as Ot,_ as Mt,G as Lt,t as ye,a as ve,v as d,d as Se,p as Pt,l as Te,q as Ft,I as we,L as Ht}from"./index-SKn09NMF.js";import{F as Rt}from"./FieldsQueryParam-NXRpAlxi.js";function mt(a,e,t){const l=a.slice();return l[10]=e[t],l}function bt(a,e,t){const l=a.slice();return l[10]=e[t],l}function _t(a,e,t){const l=a.slice();return l[15]=e[t],l}function kt(a){let e;return{c(){e=s("p"),e.innerHTML="Requires superuser <code>Authorization:TOKEN</code> header",w(e,"class","txt-hint txt-sm txt-right")},m(t,l){r(t,e,l)},d(t){t&&d(e)}}}function ht(a){let e,t,l,u,c,f,b,m,$,h,g,A,T,O,R,M,I,J,S,Q,L,q,k,P,te,Y,U,re,G,K,X;function fe(y,C){var V,W,H;return C&1&&(f=null),f==null&&(f=!!((H=(W=(V=y[0])==null?void 0:V.fields)==null?void 0:W.find(Kt))!=null&&H.required)),f?At:jt}let le=fe(a,-1),E=le(a);function Z(y,C){var V,W,H;return C&1&&(I=null),I==null&&(I=!!((H=(W=(V=y[0])==null?void 0:V.fields)==null?void 0:W.find(Gt))!=null&&H.required)),I?Vt:Bt}let x=Z(a,-1),F=x(a);return{c(){e=s("tr"),e.innerHTML='<td colspan="3" class="txt-hint txt-bold">Auth specific fields</td>',t=p(),l=s("tr"),u=s("td"),c=s("div"),E.c(),b=p(),m=s("span"),m.textContent="email",$=p(),h=s("td"),h.innerHTML='<span class="label">String</span>',g=p(),A=s("td"),A.textContent="Auth record email address.",T=p(),O=s("tr"),R=s("td"),M=s("div"),F.c(),J=p(),S=s("span"),S.textContent="emailVisibility",Q=p(),L=s("td"),L.innerHTML='<span class="label">Boolean</span>',q=p(),k=s("td"),k.textContent="Whether to show/hide the auth record email when fetching the record data.",P=p(),te=s("tr"),te.innerHTML='<td><div class="inline-flex"><span class="label label-success">Required</span> <span>password</span></div></td> <td><span class="label">String</span></td> <td>Auth record password.</td>',Y=p(),U=s("tr"),U.innerHTML='<td><div class="inline-flex"><span class="label label-success">Required</span> <span>passwordConfirm</span></div></td> <td><span class="label">String</span></td> <td>Auth record password confirmation.</td>',re=p(),G=s("tr"),G.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>verified</span></div></td> <td><span class="label">Boolean</span></td> <td>Indicates whether the auth record is verified or not.
                    <br/>
                    This field can be set only by superusers or auth records with &quot;Manage&quot; access.</td>`,K=p(),X=s("tr"),X.innerHTML='<td colspan="3" class="txt-hint txt-bold">Other fields</td>',w(c,"class","inline-flex"),w(M,"class","inline-flex")},m(y,C){r(y,e,C),r(y,t,C),r(y,l,C),i(l,u),i(u,c),E.m(c,null),i(c,b),i(c,m),i(l,$),i(l,h),i(l,g),i(l,A),r(y,T,C),r(y,O,C),i(O,R),i(R,M),F.m(M,null),i(M,J),i(M,S),i(O,Q),i(O,L),i(O,q),i(O,k),r(y,P,C),r(y,te,C),r(y,Y,C),r(y,U,C),r(y,re,C),r(y,G,C),r(y,K,C),r(y,X,C)},p(y,C){le!==(le=fe(y,C))&&(E.d(1),E=le(y),E&&(E.c(),E.m(c,b))),x!==(x=Z(y,C))&&(F.d(1),F=x(y),F&&(F.c(),F.m(M,J)))},d(y){y&&(d(e),d(t),d(l),d(T),d(O),d(P),d(te),d(Y),d(U),d(re),d(G),d(K),d(X)),E.d(),F.d()}}}function jt(a){let e;return{c(){e=s("span"),e.textContent="Optional",w(e,"class","label label-warning")},m(t,l){r(t,e,l)},d(t){t&&d(e)}}}function At(a){let e;return{c(){e=s("span"),e.textContent="Required",w(e,"class","label label-success")},m(t,l){r(t,e,l)},d(t){t&&d(e)}}}function Bt(a){let e;return{c(){e=s("span"),e.textContent="Optional",w(e,"class","label label-warning")},m(t,l){r(t,e,l)},d(t){t&&d(e)}}}function Vt(a){let e;return{c(){e=s("span"),e.textContent="Required",w(e,"class","label label-success")},m(t,l){r(t,e,l)},d(t){t&&d(e)}}}function Nt(a){let e;return{c(){e=s("span"),e.textContent="Required",w(e,"class","label label-success")},m(t,l){r(t,e,l)},d(t){t&&d(e)}}}function Dt(a){let e;return{c(){e=s("span"),e.textContent="Optional",w(e,"class","label label-warning")},m(t,l){r(t,e,l)},d(t){t&&d(e)}}}function Jt(a){let e,t=a[15].maxSelect===1?"id":"ids",l,u;return{c(){e=_("Relation record "),l=_(t),u=_(".")},m(c,f){r(c,e,f),r(c,l,f),r(c,u,f)},p(c,f){f&64&&t!==(t=c[15].maxSelect===1?"id":"ids")&&oe(l,t)},d(c){c&&(d(e),d(l),d(u))}}}function Et(a){let e,t,l,u,c,f,b,m,$;return{c(){e=_("File object."),t=s("br"),l=_(`
                        Set to empty value (`),u=s("code"),u.textContent="null",c=_(", "),f=s("code"),f.textContent='""',b=_(" or "),m=s("code"),m.textContent="[]",$=_(`) to delete
                        already uploaded file(s).`)},m(h,g){r(h,e,g),r(h,t,g),r(h,l,g),r(h,u,g),r(h,c,g),r(h,f,g),r(h,b,g),r(h,m,g),r(h,$,g)},p:we,d(h){h&&(d(e),d(t),d(l),d(u),d(c),d(f),d(b),d(m),d($))}}}function It(a){let e;return{c(){e=_("URL address.")},m(t,l){r(t,e,l)},p:we,d(t){t&&d(e)}}}function Ut(a){let e;return{c(){e=_("Email address.")},m(t,l){r(t,e,l)},p:we,d(t){t&&d(e)}}}function zt(a){let e;return{c(){e=_("JSON array or object.")},m(t,l){r(t,e,l)},p:we,d(t){t&&d(e)}}}function Qt(a){let e;return{c(){e=_("Number value.")},m(t,l){r(t,e,l)},p:we,d(t){t&&d(e)}}}function Wt(a){let e,t,l=a[15].autogeneratePattern&&yt();return{c(){e=_(`Plain text value.
                        `),l&&l.c(),t=Ht()},m(u,c){r(u,e,c),l&&l.m(u,c),r(u,t,c)},p(u,c){u[15].autogeneratePattern?l||(l=yt(),l.c(),l.m(t.parentNode,t)):l&&(l.d(1),l=null)},d(u){u&&(d(e),d(t)),l&&l.d(u)}}}function yt(a){let e;return{c(){e=_("It is autogenerated if not set.")},m(t,l){r(t,e,l)},d(t){t&&d(e)}}}function vt(a,e){let t,l,u,c,f,b=e[15].name+"",m,$,h,g,A=ee.getFieldValueType(e[15])+"",T,O,R,M;function I(k,P){return!k[15].required||k[15].type=="text"&&k[15].autogeneratePattern?Dt:Nt}let J=I(e),S=J(e);function Q(k,P){if(k[15].type==="text")return Wt;if(k[15].type==="number")return Qt;if(k[15].type==="json")return zt;if(k[15].type==="email")return Ut;if(k[15].type==="url")return It;if(k[15].type==="file")return Et;if(k[15].type==="relation")return Jt}let L=Q(e),q=L&&L(e);return{key:a,first:null,c(){t=s("tr"),l=s("td"),u=s("div"),S.c(),c=p(),f=s("span"),m=_(b),$=p(),h=s("td"),g=s("span"),T=_(A),O=p(),R=s("td"),q&&q.c(),M=p(),w(u,"class","inline-flex"),w(g,"class","label"),this.first=t},m(k,P){r(k,t,P),i(t,l),i(l,u),S.m(u,null),i(u,c),i(u,f),i(f,m),i(t,$),i(t,h),i(h,g),i(g,T),i(t,O),i(t,R),q&&q.m(R,null),i(t,M)},p(k,P){e=k,J!==(J=I(e))&&(S.d(1),S=J(e),S&&(S.c(),S.m(u,c))),P&64&&b!==(b=e[15].name+"")&&oe(m,b),P&64&&A!==(A=ee.getFieldValueType(e[15])+"")&&oe(T,A),L===(L=Q(e))&&q?q.p(e,P):(q&&q.d(1),q=L&&L(e),q&&(q.c(),q.m(R,null)))},d(k){k&&d(t),S.d(),q&&q.d()}}}function wt(a,e){let t,l=e[10].code+"",u,c,f,b;function m(){return e[9](e[10])}return{key:a,first:null,c(){t=s("button"),u=_(l),c=p(),w(t,"class","tab-item"),Te(t,"active",e[2]===e[10].code),this.first=t},m($,h){r($,t,h),i(t,u),i(t,c),f||(b=Ft(t,"click",m),f=!0)},p($,h){e=$,h&8&&l!==(l=e[10].code+"")&&oe(u,l),h&12&&Te(t,"active",e[2]===e[10].code)},d($){$&&d(t),f=!1,b()}}}function gt(a,e){let t,l,u,c;return l=new Ct({props:{content:e[10].body}}),{key:a,first:null,c(){t=s("div"),$e(l.$$.fragment),u=p(),w(t,"class","tab-item"),Te(t,"active",e[2]===e[10].code),this.first=t},m(f,b){r(f,t,b),qe(l,t,null),i(t,u),c=!0},p(f,b){e=f;const m={};b&8&&(m.content=e[10].body),l.$set(m),(!c||b&12)&&Te(t,"active",e[2]===e[10].code)},i(f){c||(ye(l.$$.fragment,f),c=!0)},o(f){ve(l.$$.fragment,f),c=!1},d(f){f&&d(t),Se(l)}}}function Yt(a){var at,st,ot,rt;let e,t,l=a[0].name+"",u,c,f,b,m,$,h,g=a[0].name+"",A,T,O,R,M,I,J,S,Q,L,q,k,P,te,Y,U,re,G,K=a[0].name+"",X,fe,le,E,Z,x,F,y,C,V,W,H=[],Ne=new Map,Oe,pe,Me,ne,Le,De,me,ie,Pe,Je,Fe,Ee,j,Ie,de,Ue,ze,Qe,He,We,Re,Ye,Ge,Ke,je,Xe,Ze,ce,Ae,be,Be,ae,_e,z=[],xe=new Map,et,ke,N=[],tt=new Map,se;S=new Tt({props:{js:`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${a[5]}');

...

// example create data
const data = ${JSON.stringify(Object.assign({},a[4],ee.dummyCollectionSchemaData(a[0],!0)),null,4)};

const record = await pb.collection('${(at=a[0])==null?void 0:at.name}').create(data);
`+(a[1]?`
// (optional) send an email verification request
await pb.collection('${(st=a[0])==null?void 0:st.name}').requestVerification('test@example.com');
`:""),dart:`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${a[5]}');

...

// example create body
final body = <String, dynamic>${JSON.stringify(Object.assign({},a[4],ee.dummyCollectionSchemaData(a[0],!0)),null,2)};

final record = await pb.collection('${(ot=a[0])==null?void 0:ot.name}').create(body: body);
`+(a[1]?`
// (optional) send an email verification request
await pb.collection('${(rt=a[0])==null?void 0:rt.name}').requestVerification('test@example.com');
`:"")}});let D=a[7]&&kt(),B=a[1]&&ht(a),ge=ue(a[6]);const lt=n=>n[15].name;for(let n=0;n<ge.length;n+=1){let o=_t(a,ge,n),v=lt(o);Ne.set(v,H[n]=vt(v,o))}de=new Ct({props:{content:"?expand=relField1,relField2.subRelField"}}),ce=new Rt({});let Ce=ue(a[3]);const nt=n=>n[10].code;for(let n=0;n<Ce.length;n+=1){let o=bt(a,Ce,n),v=nt(o);xe.set(v,z[n]=wt(v,o))}let he=ue(a[3]);const it=n=>n[10].code;for(let n=0;n<he.length;n+=1){let o=mt(a,he,n),v=it(o);tt.set(v,N[n]=gt(v,o))}return{c(){e=s("h3"),t=_("Create ("),u=_(l),c=_(")"),f=p(),b=s("div"),m=s("p"),$=_("Create a new "),h=s("strong"),A=_(g),T=_(" record."),O=p(),R=s("p"),R.innerHTML=`Body parameters could be sent as <code>application/json</code> or
        <code>multipart/form-data</code>.`,M=p(),I=s("p"),I.innerHTML=`File upload is supported only via <code>multipart/form-data</code>.
        <br/>
        For more info and examples you could check the detailed
        <a href="https://pocketbase.io/docs/files-handling" target="_blank" rel="noopener noreferrer">Files upload and handling docs
        </a>.`,J=p(),$e(S.$$.fragment),Q=p(),L=s("h6"),L.textContent="API details",q=p(),k=s("div"),P=s("strong"),P.textContent="POST",te=p(),Y=s("div"),U=s("p"),re=_("/api/collections/"),G=s("strong"),X=_(K),fe=_("/records"),le=p(),D&&D.c(),E=p(),Z=s("div"),Z.textContent="Body Parameters",x=p(),F=s("table"),y=s("thead"),y.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr>',C=p(),V=s("tbody"),B&&B.c(),W=p();for(let n=0;n<H.length;n+=1)H[n].c();Oe=p(),pe=s("div"),pe.textContent="Query parameters",Me=p(),ne=s("table"),Le=s("thead"),Le.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr>',De=p(),me=s("tbody"),ie=s("tr"),Pe=s("td"),Pe.textContent="expand",Je=p(),Fe=s("td"),Fe.innerHTML='<span class="label">String</span>',Ee=p(),j=s("td"),Ie=_(`Auto expand relations when returning the created record. Ex.:
                `),$e(de.$$.fragment),Ue=_(`
                Supports up to 6-levels depth nested relations expansion. `),ze=s("br"),Qe=_(`
                The expanded relations will be appended to the record under the
                `),He=s("code"),He.textContent="expand",We=_(" property (eg. "),Re=s("code"),Re.textContent='"expand": {"relField1": {...}, ...}',Ye=_(`).
                `),Ge=s("br"),Ke=_(`
                Only the relations to which the request user has permissions to `),je=s("strong"),je.textContent="view",Xe=_(" will be expanded."),Ze=p(),$e(ce.$$.fragment),Ae=p(),be=s("div"),be.textContent="Responses",Be=p(),ae=s("div"),_e=s("div");for(let n=0;n<z.length;n+=1)z[n].c();et=p(),ke=s("div");for(let n=0;n<N.length;n+=1)N[n].c();w(e,"class","m-b-sm"),w(b,"class","content txt-lg m-b-sm"),w(L,"class","m-b-xs"),w(P,"class","label label-primary"),w(Y,"class","content"),w(k,"class","alert alert-success"),w(Z,"class","section-title"),w(F,"class","table-compact table-border m-b-base"),w(pe,"class","section-title"),w(ne,"class","table-compact table-border m-b-base"),w(be,"class","section-title"),w(_e,"class","tabs-header compact combined left"),w(ke,"class","tabs-content"),w(ae,"class","tabs")},m(n,o){r(n,e,o),i(e,t),i(e,u),i(e,c),r(n,f,o),r(n,b,o),i(b,m),i(m,$),i(m,h),i(h,A),i(m,T),i(b,O),i(b,R),i(b,M),i(b,I),r(n,J,o),qe(S,n,o),r(n,Q,o),r(n,L,o),r(n,q,o),r(n,k,o),i(k,P),i(k,te),i(k,Y),i(Y,U),i(U,re),i(U,G),i(G,X),i(U,fe),i(k,le),D&&D.m(k,null),r(n,E,o),r(n,Z,o),r(n,x,o),r(n,F,o),i(F,y),i(F,C),i(F,V),B&&B.m(V,null),i(V,W);for(let v=0;v<H.length;v+=1)H[v]&&H[v].m(V,null);r(n,Oe,o),r(n,pe,o),r(n,Me,o),r(n,ne,o),i(ne,Le),i(ne,De),i(ne,me),i(me,ie),i(ie,Pe),i(ie,Je),i(ie,Fe),i(ie,Ee),i(ie,j),i(j,Ie),qe(de,j,null),i(j,Ue),i(j,ze),i(j,Qe),i(j,He),i(j,We),i(j,Re),i(j,Ye),i(j,Ge),i(j,Ke),i(j,je),i(j,Xe),i(me,Ze),qe(ce,me,null),r(n,Ae,o),r(n,be,o),r(n,Be,o),r(n,ae,o),i(ae,_e);for(let v=0;v<z.length;v+=1)z[v]&&z[v].m(_e,null);i(ae,et),i(ae,ke);for(let v=0;v<N.length;v+=1)N[v]&&N[v].m(ke,null);se=!0},p(n,[o]){var dt,ct,ut,ft;(!se||o&1)&&l!==(l=n[0].name+"")&&oe(u,l),(!se||o&1)&&g!==(g=n[0].name+"")&&oe(A,g);const v={};o&51&&(v.js=`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${n[5]}');

...

// example create data
const data = ${JSON.stringify(Object.assign({},n[4],ee.dummyCollectionSchemaData(n[0],!0)),null,4)};

const record = await pb.collection('${(dt=n[0])==null?void 0:dt.name}').create(data);
`+(n[1]?`
// (optional) send an email verification request
await pb.collection('${(ct=n[0])==null?void 0:ct.name}').requestVerification('test@example.com');
`:"")),o&51&&(v.dart=`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${n[5]}');

...

// example create body
final body = <String, dynamic>${JSON.stringify(Object.assign({},n[4],ee.dummyCollectionSchemaData(n[0],!0)),null,2)};

final record = await pb.collection('${(ut=n[0])==null?void 0:ut.name}').create(body: body);
`+(n[1]?`
// (optional) send an email verification request
await pb.collection('${(ft=n[0])==null?void 0:ft.name}').requestVerification('test@example.com');
`:"")),S.$set(v),(!se||o&1)&&K!==(K=n[0].name+"")&&oe(X,K),n[7]?D||(D=kt(),D.c(),D.m(k,null)):D&&(D.d(1),D=null),n[1]?B?B.p(n,o):(B=ht(n),B.c(),B.m(V,W)):B&&(B.d(1),B=null),o&64&&(ge=ue(n[6]),H=Ve(H,o,lt,1,n,ge,Ne,V,pt,vt,null,_t)),o&12&&(Ce=ue(n[3]),z=Ve(z,o,nt,1,n,Ce,xe,_e,pt,wt,null,bt)),o&12&&(he=ue(n[3]),Ot(),N=Ve(N,o,it,1,n,he,tt,ke,Mt,gt,null,mt),Lt())},i(n){if(!se){ye(S.$$.fragment,n),ye(de.$$.fragment,n),ye(ce.$$.fragment,n);for(let o=0;o<he.length;o+=1)ye(N[o]);se=!0}},o(n){ve(S.$$.fragment,n),ve(de.$$.fragment,n),ve(ce.$$.fragment,n);for(let o=0;o<N.length;o+=1)ve(N[o]);se=!1},d(n){n&&(d(e),d(f),d(b),d(J),d(Q),d(L),d(q),d(k),d(E),d(Z),d(x),d(F),d(Oe),d(pe),d(Me),d(ne),d(Ae),d(be),d(Be),d(ae)),Se(S,n),D&&D.d(),B&&B.d();for(let o=0;o<H.length;o+=1)H[o].d();Se(de),Se(ce);for(let o=0;o<z.length;o+=1)z[o].d();for(let o=0;o<N.length;o+=1)N[o].d()}}}const Gt=a=>a.name=="emailVisibility",Kt=a=>a.name=="email";function Xt(a,e,t){let l,u,c,f,b,{collection:m}=e,$=200,h=[],g={};const A=T=>t(2,$=T.code);return a.$$set=T=>{"collection"in T&&t(0,m=T.collection)},a.$$.update=()=>{var T,O,R;a.$$.dirty&1&&t(1,l=m.type==="auth"),a.$$.dirty&1&&t(7,u=(m==null?void 0:m.createRule)===null),a.$$.dirty&2&&t(8,c=l?["password","verified","email","emailVisibility"]:[]),a.$$.dirty&257&&t(6,f=((T=m==null?void 0:m.fields)==null?void 0:T.filter(M=>!M.hidden&&M.type!="autodate"&&!c.includes(M.name)))||[]),a.$$.dirty&1&&t(3,h=[{code:200,body:JSON.stringify(ee.dummyCollectionRecord(m),null,2)},{code:400,body:`
                {
                  "code": 400,
                  "message": "Failed to create record.",
                  "data": {
                    "${(R=(O=m==null?void 0:m.fields)==null?void 0:O[0])==null?void 0:R.name}": {
                      "code": "validation_required",
                      "message": "Missing required value."
                    }
                  }
                }
            `},{code:403,body:`
                {
                  "code": 403,
                  "message": "You are not allowed to perform this request.",
                  "data": {}
                }
            `}]),a.$$.dirty&2&&(l?t(4,g={password:"12345678",passwordConfirm:"12345678"}):t(4,g={}))},t(5,b=ee.getApiExampleUrl(Pt.baseURL)),[m,l,$,h,g,b,f,u,c,A]}class el extends $t{constructor(e){super(),qt(this,e,Xt,Yt,St,{collection:0})}}export{el as default};
