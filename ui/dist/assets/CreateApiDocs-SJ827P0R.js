import{S as $t,i as qt,s as St,V as Tt,J as ee,X as ue,W as Ct,h as r,d as $e,t as ye,a as ve,I as oe,Z as Ve,_ as pt,C as Ot,$ as Mt,D as Lt,l as d,n as i,m as qe,u as a,A as _,v as p,c as Se,w,p as Pt,k as Te,o as Ft,L as Ht,H as we}from"./index-CVqXRFk8.js";import{F as Rt}from"./FieldsQueryParam-B6maRWKB.js";function mt(s,e,t){const l=s.slice();return l[10]=e[t],l}function bt(s,e,t){const l=s.slice();return l[10]=e[t],l}function _t(s,e,t){const l=s.slice();return l[15]=e[t],l}function kt(s){let e;return{c(){e=a("p"),e.innerHTML="Requires superuser <code>Authorization:TOKEN</code> header",w(e,"class","txt-hint txt-sm txt-right")},m(t,l){d(t,e,l)},d(t){t&&r(e)}}}function ht(s){let e,t,l,u,c,f,b,m,$,h,g,B,T,O,R,M,I,J,S,W,L,q,k,P,te,K,U,re,X,Y,Z;function fe(y,C){var V,z,H;return C&1&&(f=null),f==null&&(f=!!((H=(z=(V=y[0])==null?void 0:V.fields)==null?void 0:z.find(Yt))!=null&&H.required)),f?Bt:At}let le=fe(s,-1),E=le(s);function G(y,C){var V,z,H;return C&1&&(I=null),I==null&&(I=!!((H=(z=(V=y[0])==null?void 0:V.fields)==null?void 0:z.find(Xt))!=null&&H.required)),I?Vt:jt}let x=G(s,-1),F=x(s);return{c(){e=a("tr"),e.innerHTML='<td colspan="3" class="txt-hint txt-bold">Auth specific fields</td>',t=p(),l=a("tr"),u=a("td"),c=a("div"),E.c(),b=p(),m=a("span"),m.textContent="email",$=p(),h=a("td"),h.innerHTML='<span class="label">String</span>',g=p(),B=a("td"),B.textContent="Auth record email address.",T=p(),O=a("tr"),R=a("td"),M=a("div"),F.c(),J=p(),S=a("span"),S.textContent="emailVisibility",W=p(),L=a("td"),L.innerHTML='<span class="label">Boolean</span>',q=p(),k=a("td"),k.textContent="Whether to show/hide the auth record email when fetching the record data.",P=p(),te=a("tr"),te.innerHTML='<td><div class="inline-flex"><span class="label label-success">Required</span> <span>password</span></div></td> <td><span class="label">String</span></td> <td>Auth record password.</td>',K=p(),U=a("tr"),U.innerHTML='<td><div class="inline-flex"><span class="label label-success">Required</span> <span>passwordConfirm</span></div></td> <td><span class="label">String</span></td> <td>Auth record password confirmation.</td>',re=p(),X=a("tr"),X.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>verified</span></div></td> <td><span class="label">Boolean</span></td> <td>Indicates whether the auth record is verified or not.
                    <br/>
                    This field can be set only by superusers or auth records with &quot;Manage&quot; access.</td>`,Y=p(),Z=a("tr"),Z.innerHTML='<td colspan="3" class="txt-hint txt-bold">Other fields</td>',w(c,"class","inline-flex"),w(M,"class","inline-flex")},m(y,C){d(y,e,C),d(y,t,C),d(y,l,C),i(l,u),i(u,c),E.m(c,null),i(c,b),i(c,m),i(l,$),i(l,h),i(l,g),i(l,B),d(y,T,C),d(y,O,C),i(O,R),i(R,M),F.m(M,null),i(M,J),i(M,S),i(O,W),i(O,L),i(O,q),i(O,k),d(y,P,C),d(y,te,C),d(y,K,C),d(y,U,C),d(y,re,C),d(y,X,C),d(y,Y,C),d(y,Z,C)},p(y,C){le!==(le=fe(y,C))&&(E.d(1),E=le(y),E&&(E.c(),E.m(c,b))),x!==(x=G(y,C))&&(F.d(1),F=x(y),F&&(F.c(),F.m(M,J)))},d(y){y&&(r(e),r(t),r(l),r(T),r(O),r(P),r(te),r(K),r(U),r(re),r(X),r(Y),r(Z)),E.d(),F.d()}}}function At(s){let e;return{c(){e=a("span"),e.textContent="Optional",w(e,"class","label label-warning")},m(t,l){d(t,e,l)},d(t){t&&r(e)}}}function Bt(s){let e;return{c(){e=a("span"),e.textContent="Required",w(e,"class","label label-success")},m(t,l){d(t,e,l)},d(t){t&&r(e)}}}function jt(s){let e;return{c(){e=a("span"),e.textContent="Optional",w(e,"class","label label-warning")},m(t,l){d(t,e,l)},d(t){t&&r(e)}}}function Vt(s){let e;return{c(){e=a("span"),e.textContent="Required",w(e,"class","label label-success")},m(t,l){d(t,e,l)},d(t){t&&r(e)}}}function Dt(s){let e;return{c(){e=a("span"),e.textContent="Required",w(e,"class","label label-success")},m(t,l){d(t,e,l)},d(t){t&&r(e)}}}function Nt(s){let e;return{c(){e=a("span"),e.textContent="Optional",w(e,"class","label label-warning")},m(t,l){d(t,e,l)},d(t){t&&r(e)}}}function Jt(s){let e,t=s[15].maxSelect===1?"id":"ids",l,u;return{c(){e=_("Relation record "),l=_(t),u=_(".")},m(c,f){d(c,e,f),d(c,l,f),d(c,u,f)},p(c,f){f&64&&t!==(t=c[15].maxSelect===1?"id":"ids")&&oe(l,t)},d(c){c&&(r(e),r(l),r(u))}}}function Et(s){let e,t,l,u,c,f,b,m,$;return{c(){e=_("File object."),t=a("br"),l=_(`
                        Set to empty value (`),u=a("code"),u.textContent="null",c=_(", "),f=a("code"),f.textContent='""',b=_(" or "),m=a("code"),m.textContent="[]",$=_(`) to delete
                        already uploaded file(s).`)},m(h,g){d(h,e,g),d(h,t,g),d(h,l,g),d(h,u,g),d(h,c,g),d(h,f,g),d(h,b,g),d(h,m,g),d(h,$,g)},p:we,d(h){h&&(r(e),r(t),r(l),r(u),r(c),r(f),r(b),r(m),r($))}}}function It(s){let e;return{c(){e=_("URL address.")},m(t,l){d(t,e,l)},p:we,d(t){t&&r(e)}}}function Ut(s){let e;return{c(){e=_("Email address.")},m(t,l){d(t,e,l)},p:we,d(t){t&&r(e)}}}function Qt(s){let e;return{c(){e=_("JSON array or object.")},m(t,l){d(t,e,l)},p:we,d(t){t&&r(e)}}}function Wt(s){let e;return{c(){e=_("Number value.")},m(t,l){d(t,e,l)},p:we,d(t){t&&r(e)}}}function zt(s){let e,t,l=s[15].autogeneratePattern&&yt();return{c(){e=_(`Plain text value.
                        `),l&&l.c(),t=Ht()},m(u,c){d(u,e,c),l&&l.m(u,c),d(u,t,c)},p(u,c){u[15].autogeneratePattern?l||(l=yt(),l.c(),l.m(t.parentNode,t)):l&&(l.d(1),l=null)},d(u){u&&(r(e),r(t)),l&&l.d(u)}}}function yt(s){let e;return{c(){e=_("It is autogenerated if not set.")},m(t,l){d(t,e,l)},d(t){t&&r(e)}}}function vt(s,e){let t,l,u,c,f,b=e[15].name+"",m,$,h,g,B=ee.getFieldValueType(e[15])+"",T,O,R,M;function I(k,P){return!k[15].required||k[15].type=="text"&&k[15].autogeneratePattern?Nt:Dt}let J=I(e),S=J(e);function W(k,P){if(k[15].type==="text")return zt;if(k[15].type==="number")return Wt;if(k[15].type==="json")return Qt;if(k[15].type==="email")return Ut;if(k[15].type==="url")return It;if(k[15].type==="file")return Et;if(k[15].type==="relation")return Jt}let L=W(e),q=L&&L(e);return{key:s,first:null,c(){t=a("tr"),l=a("td"),u=a("div"),S.c(),c=p(),f=a("span"),m=_(b),$=p(),h=a("td"),g=a("span"),T=_(B),O=p(),R=a("td"),q&&q.c(),M=p(),w(u,"class","inline-flex"),w(g,"class","label"),this.first=t},m(k,P){d(k,t,P),i(t,l),i(l,u),S.m(u,null),i(u,c),i(u,f),i(f,m),i(t,$),i(t,h),i(h,g),i(g,T),i(t,O),i(t,R),q&&q.m(R,null),i(t,M)},p(k,P){e=k,J!==(J=I(e))&&(S.d(1),S=J(e),S&&(S.c(),S.m(u,c))),P&64&&b!==(b=e[15].name+"")&&oe(m,b),P&64&&B!==(B=ee.getFieldValueType(e[15])+"")&&oe(T,B),L===(L=W(e))&&q?q.p(e,P):(q&&q.d(1),q=L&&L(e),q&&(q.c(),q.m(R,null)))},d(k){k&&r(t),S.d(),q&&q.d()}}}function wt(s,e){let t,l=e[10].code+"",u,c,f,b;function m(){return e[9](e[10])}return{key:s,first:null,c(){t=a("button"),u=_(l),c=p(),w(t,"class","tab-item"),Te(t,"active",e[2]===e[10].code),this.first=t},m($,h){d($,t,h),i(t,u),i(t,c),f||(b=Ft(t,"click",m),f=!0)},p($,h){e=$,h&8&&l!==(l=e[10].code+"")&&oe(u,l),h&12&&Te(t,"active",e[2]===e[10].code)},d($){$&&r(t),f=!1,b()}}}function gt(s,e){let t,l,u,c;return l=new Ct({props:{content:e[10].body}}),{key:s,first:null,c(){t=a("div"),Se(l.$$.fragment),u=p(),w(t,"class","tab-item"),Te(t,"active",e[2]===e[10].code),this.first=t},m(f,b){d(f,t,b),qe(l,t,null),i(t,u),c=!0},p(f,b){e=f;const m={};b&8&&(m.content=e[10].body),l.$set(m),(!c||b&12)&&Te(t,"active",e[2]===e[10].code)},i(f){c||(ve(l.$$.fragment,f),c=!0)},o(f){ye(l.$$.fragment,f),c=!1},d(f){f&&r(t),$e(l)}}}function Kt(s){var st,at,ot,rt;let e,t,l=s[0].name+"",u,c,f,b,m,$,h,g=s[0].name+"",B,T,O,R,M,I,J,S,W,L,q,k,P,te,K,U,re,X,Y=s[0].name+"",Z,fe,le,E,G,x,F,y,C,V,z,H=[],De=new Map,Oe,pe,Me,ne,Le,Ne,me,ie,Pe,Je,Fe,Ee,A,Ie,de,Ue,Qe,We,He,ze,Re,Ke,Xe,Ye,Ae,Ze,Ge,ce,Be,be,je,se,_e,Q=[],xe=new Map,et,ke,D=[],tt=new Map,ae;S=new Tt({props:{js:`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${s[5]}');

...

// example create data
const data = ${JSON.stringify(Object.assign({},s[4],ee.dummyCollectionSchemaData(s[0],!0)),null,4)};

const record = await pb.collection('${(st=s[0])==null?void 0:st.name}').create(data);
`+(s[1]?`
// (optional) send an email verification request
await pb.collection('${(at=s[0])==null?void 0:at.name}').requestVerification('test@example.com');
`:""),dart:`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${s[5]}');

...

// example create body
final body = <String, dynamic>${JSON.stringify(Object.assign({},s[4],ee.dummyCollectionSchemaData(s[0],!0)),null,2)};

final record = await pb.collection('${(ot=s[0])==null?void 0:ot.name}').create(body: body);
`+(s[1]?`
// (optional) send an email verification request
await pb.collection('${(rt=s[0])==null?void 0:rt.name}').requestVerification('test@example.com');
`:"")}});let N=s[7]&&kt(),j=s[1]&&ht(s),ge=ue(s[6]);const lt=n=>n[15].name;for(let n=0;n<ge.length;n+=1){let o=_t(s,ge,n),v=lt(o);De.set(v,H[n]=vt(v,o))}de=new Ct({props:{content:"?expand=relField1,relField2.subRelField"}}),ce=new Rt({});let Ce=ue(s[3]);const nt=n=>n[10].code;for(let n=0;n<Ce.length;n+=1){let o=bt(s,Ce,n),v=nt(o);xe.set(v,Q[n]=wt(v,o))}let he=ue(s[3]);const it=n=>n[10].code;for(let n=0;n<he.length;n+=1){let o=mt(s,he,n),v=it(o);tt.set(v,D[n]=gt(v,o))}return{c(){e=a("h3"),t=_("Create ("),u=_(l),c=_(")"),f=p(),b=a("div"),m=a("p"),$=_("Create a new "),h=a("strong"),B=_(g),T=_(" record."),O=p(),R=a("p"),R.innerHTML=`Body parameters could be sent as <code>application/json</code> or
        <code>multipart/form-data</code>.`,M=p(),I=a("p"),I.innerHTML=`File upload is supported only via <code>multipart/form-data</code>.
        <br/>
        For more info and examples you could check the detailed
        <a href="https://pocketbase.io/docs/files-handling" target="_blank" rel="noopener noreferrer">Files upload and handling docs
        </a>.`,J=p(),Se(S.$$.fragment),W=p(),L=a("h6"),L.textContent="API details",q=p(),k=a("div"),P=a("strong"),P.textContent="POST",te=p(),K=a("div"),U=a("p"),re=_("/api/collections/"),X=a("strong"),Z=_(Y),fe=_("/records"),le=p(),N&&N.c(),E=p(),G=a("div"),G.textContent="Body Parameters",x=p(),F=a("table"),y=a("thead"),y.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr>',C=p(),V=a("tbody"),j&&j.c(),z=p();for(let n=0;n<H.length;n+=1)H[n].c();Oe=p(),pe=a("div"),pe.textContent="Query parameters",Me=p(),ne=a("table"),Le=a("thead"),Le.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr>',Ne=p(),me=a("tbody"),ie=a("tr"),Pe=a("td"),Pe.textContent="expand",Je=p(),Fe=a("td"),Fe.innerHTML='<span class="label">String</span>',Ee=p(),A=a("td"),Ie=_(`Auto expand relations when returning the created record. Ex.:
                `),Se(de.$$.fragment),Ue=_(`
                Supports up to 6-levels depth nested relations expansion. `),Qe=a("br"),We=_(`
                The expanded relations will be appended to the record under the
                `),He=a("code"),He.textContent="expand",ze=_(" property (eg. "),Re=a("code"),Re.textContent='"expand": {"relField1": {...}, ...}',Ke=_(`).
                `),Xe=a("br"),Ye=_(`
                Only the relations to which the request user has permissions to `),Ae=a("strong"),Ae.textContent="view",Ze=_(" will be expanded."),Ge=p(),Se(ce.$$.fragment),Be=p(),be=a("div"),be.textContent="Responses",je=p(),se=a("div"),_e=a("div");for(let n=0;n<Q.length;n+=1)Q[n].c();et=p(),ke=a("div");for(let n=0;n<D.length;n+=1)D[n].c();w(e,"class","m-b-sm"),w(b,"class","content txt-lg m-b-sm"),w(L,"class","m-b-xs"),w(P,"class","label label-primary"),w(K,"class","content"),w(k,"class","alert alert-success"),w(G,"class","section-title"),w(F,"class","table-compact table-border m-b-base"),w(pe,"class","section-title"),w(ne,"class","table-compact table-border m-b-base"),w(be,"class","section-title"),w(_e,"class","tabs-header compact combined left"),w(ke,"class","tabs-content"),w(se,"class","tabs")},m(n,o){d(n,e,o),i(e,t),i(e,u),i(e,c),d(n,f,o),d(n,b,o),i(b,m),i(m,$),i(m,h),i(h,B),i(m,T),i(b,O),i(b,R),i(b,M),i(b,I),d(n,J,o),qe(S,n,o),d(n,W,o),d(n,L,o),d(n,q,o),d(n,k,o),i(k,P),i(k,te),i(k,K),i(K,U),i(U,re),i(U,X),i(X,Z),i(U,fe),i(k,le),N&&N.m(k,null),d(n,E,o),d(n,G,o),d(n,x,o),d(n,F,o),i(F,y),i(F,C),i(F,V),j&&j.m(V,null),i(V,z);for(let v=0;v<H.length;v+=1)H[v]&&H[v].m(V,null);d(n,Oe,o),d(n,pe,o),d(n,Me,o),d(n,ne,o),i(ne,Le),i(ne,Ne),i(ne,me),i(me,ie),i(ie,Pe),i(ie,Je),i(ie,Fe),i(ie,Ee),i(ie,A),i(A,Ie),qe(de,A,null),i(A,Ue),i(A,Qe),i(A,We),i(A,He),i(A,ze),i(A,Re),i(A,Ke),i(A,Xe),i(A,Ye),i(A,Ae),i(A,Ze),i(me,Ge),qe(ce,me,null),d(n,Be,o),d(n,be,o),d(n,je,o),d(n,se,o),i(se,_e);for(let v=0;v<Q.length;v+=1)Q[v]&&Q[v].m(_e,null);i(se,et),i(se,ke);for(let v=0;v<D.length;v+=1)D[v]&&D[v].m(ke,null);ae=!0},p(n,[o]){var dt,ct,ut,ft;(!ae||o&1)&&l!==(l=n[0].name+"")&&oe(u,l),(!ae||o&1)&&g!==(g=n[0].name+"")&&oe(B,g);const v={};o&51&&(v.js=`
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
`:"")),S.$set(v),(!ae||o&1)&&Y!==(Y=n[0].name+"")&&oe(Z,Y),n[7]?N||(N=kt(),N.c(),N.m(k,null)):N&&(N.d(1),N=null),n[1]?j?j.p(n,o):(j=ht(n),j.c(),j.m(V,z)):j&&(j.d(1),j=null),o&64&&(ge=ue(n[6]),H=Ve(H,o,lt,1,n,ge,De,V,pt,vt,null,_t)),o&12&&(Ce=ue(n[3]),Q=Ve(Q,o,nt,1,n,Ce,xe,_e,pt,wt,null,bt)),o&12&&(he=ue(n[3]),Ot(),D=Ve(D,o,it,1,n,he,tt,ke,Mt,gt,null,mt),Lt())},i(n){if(!ae){ve(S.$$.fragment,n),ve(de.$$.fragment,n),ve(ce.$$.fragment,n);for(let o=0;o<he.length;o+=1)ve(D[o]);ae=!0}},o(n){ye(S.$$.fragment,n),ye(de.$$.fragment,n),ye(ce.$$.fragment,n);for(let o=0;o<D.length;o+=1)ye(D[o]);ae=!1},d(n){n&&(r(e),r(f),r(b),r(J),r(W),r(L),r(q),r(k),r(E),r(G),r(x),r(F),r(Oe),r(pe),r(Me),r(ne),r(Be),r(be),r(je),r(se)),$e(S,n),N&&N.d(),j&&j.d();for(let o=0;o<H.length;o+=1)H[o].d();$e(de),$e(ce);for(let o=0;o<Q.length;o+=1)Q[o].d();for(let o=0;o<D.length;o+=1)D[o].d()}}}const Xt=s=>s.name=="emailVisibility",Yt=s=>s.name=="email";function Zt(s,e,t){let l,u,c,f,b,{collection:m}=e,$=200,h=[],g={};const B=T=>t(2,$=T.code);return s.$$set=T=>{"collection"in T&&t(0,m=T.collection)},s.$$.update=()=>{var T,O,R;s.$$.dirty&1&&t(1,l=m.type==="auth"),s.$$.dirty&1&&t(7,u=(m==null?void 0:m.createRule)===null),s.$$.dirty&2&&t(8,c=l?["password","verified","email","emailVisibility"]:[]),s.$$.dirty&257&&t(6,f=((T=m==null?void 0:m.fields)==null?void 0:T.filter(M=>!M.hidden&&M.type!="autodate"&&!c.includes(M.name)))||[]),s.$$.dirty&1&&t(3,h=[{code:200,body:JSON.stringify(ee.dummyCollectionRecord(m),null,2)},{code:400,body:`
                {
                  "status": 400,
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
                  "status": 403,
                  "message": "You are not allowed to perform this request.",
                  "data": {}
                }
            `}]),s.$$.dirty&2&&(l?t(4,g={password:"12345678",passwordConfirm:"12345678"}):t(4,g={}))},t(5,b=ee.getApiExampleUrl(Pt.baseURL)),[m,l,$,h,g,b,f,u,c,B]}class el extends $t{constructor(e){super(),qt(this,e,Zt,Kt,St,{collection:0})}}export{el as default};
