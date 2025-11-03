import{S as $t,i as Mt,s as Ot,ab as St,O as se,ac as Tt,h as d,d as ge,t as _e,a as he,I as ee,ad as Je,ae as bt,C as qt,af as Rt,D as Ht,l as o,n as a,m as we,u as s,A as _,v as f,c as Ce,w as k,J as ye,p as Pt,k as Te,o as Lt,H as te}from"./index-B4ZsHsKR.js";import{F as Dt}from"./FieldsQueryParam-K1y4zYh0.js";function mt(r,e,t){const n=r.slice();return n[10]=e[t],n}function _t(r,e,t){const n=r.slice();return n[10]=e[t],n}function ht(r,e,t){const n=r.slice();return n[15]=e[t],n}function yt(r){let e;return{c(){e=s("p"),e.innerHTML=`<em>Note that in case of a password change all previously issued tokens for the current record
                will be automatically invalidated and if you want your user to remain signed in you need to
                reauthenticate manually after the update call.</em>`},m(t,n){o(t,e,n)},d(t){t&&d(e)}}}function kt(r){let e;return{c(){e=s("p"),e.innerHTML="Requires superuser <code>Authorization:TOKEN</code> header",k(e,"class","txt-hint txt-sm txt-right")},m(t,n){o(t,e,n)},d(t){t&&d(e)}}}function vt(r){let e,t,n,b,p,c,u,m,O,T,H,P,$,M,q,L,J,j,S,R,D,v,g,w;function x(h,C){var le,z,ne;return C&1&&(m=null),m==null&&(m=!!((ne=(z=(le=h[0])==null?void 0:le.fields)==null?void 0:z.find(Kt))!=null&&ne.required)),m?Bt:Ft}let Q=x(r,-1),B=Q(r);return{c(){e=s("tr"),e.innerHTML='<td colspan="3" class="txt-hint txt-bold">Auth specific fields</td>',t=f(),n=s("tr"),n.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>email</span></div></td> <td><span class="label">String</span></td> <td>The auth record email address.
                    <br/>
                    This field can be updated only by superusers or auth records with &quot;Manage&quot; access.
                    <br/>
                    Regular accounts can update their email by calling &quot;Request email change&quot;.</td>`,b=f(),p=s("tr"),c=s("td"),u=s("div"),B.c(),O=f(),T=s("span"),T.textContent="emailVisibility",H=f(),P=s("td"),P.innerHTML='<span class="label">Boolean</span>',$=f(),M=s("td"),M.textContent="Whether to show/hide the auth record email when fetching the record data.",q=f(),L=s("tr"),L.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>oldPassword</span></div></td> <td><span class="label">String</span></td> <td>Old auth record password.
                    <br/>
                    This field is required only when changing the record password. Superusers and auth records
                    with &quot;Manage&quot; access can skip this field.</td>`,J=f(),j=s("tr"),j.innerHTML='<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>password</span></div></td> <td><span class="label">String</span></td> <td>New auth record password.</td>',S=f(),R=s("tr"),R.innerHTML='<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>passwordConfirm</span></div></td> <td><span class="label">String</span></td> <td>New auth record password confirmation.</td>',D=f(),v=s("tr"),v.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>verified</span></div></td> <td><span class="label">Boolean</span></td> <td>Indicates whether the auth record is verified or not.
                    <br/>
                    This field can be set only by superusers or auth records with &quot;Manage&quot; access.</td>`,g=f(),w=s("tr"),w.innerHTML='<td colspan="3" class="txt-hint txt-bold">Other fields</td>',k(u,"class","inline-flex")},m(h,C){o(h,e,C),o(h,t,C),o(h,n,C),o(h,b,C),o(h,p,C),a(p,c),a(c,u),B.m(u,null),a(u,O),a(u,T),a(p,H),a(p,P),a(p,$),a(p,M),o(h,q,C),o(h,L,C),o(h,J,C),o(h,j,C),o(h,S,C),o(h,R,C),o(h,D,C),o(h,v,C),o(h,g,C),o(h,w,C)},p(h,C){Q!==(Q=x(h,C))&&(B.d(1),B=Q(h),B&&(B.c(),B.m(u,O)))},d(h){h&&(d(e),d(t),d(n),d(b),d(p),d(q),d(L),d(J),d(j),d(S),d(R),d(D),d(v),d(g),d(w)),B.d()}}}function Ft(r){let e;return{c(){e=s("span"),e.textContent="Optional",k(e,"class","label label-warning")},m(t,n){o(t,e,n)},d(t){t&&d(e)}}}function Bt(r){let e;return{c(){e=s("span"),e.textContent="Required",k(e,"class","label label-success")},m(t,n){o(t,e,n)},d(t){t&&d(e)}}}function Nt(r){let e;return{c(){e=s("span"),e.textContent="Optional",k(e,"class","label label-warning")},m(t,n){o(t,e,n)},d(t){t&&d(e)}}}function At(r){let e;return{c(){e=s("span"),e.textContent="Required",k(e,"class","label label-success")},m(t,n){o(t,e,n)},d(t){t&&d(e)}}}function Et(r){let e,t=r[15].maxSelect==1?"id":"ids",n,b;return{c(){e=_("Relation record "),n=_(t),b=_(".")},m(p,c){o(p,e,c),o(p,n,c),o(p,b,c)},p(p,c){c&32&&t!==(t=p[15].maxSelect==1?"id":"ids")&&ee(n,t)},d(p){p&&(d(e),d(n),d(b))}}}function It(r){let e,t,n,b,p;return{c(){e=_("File object."),t=s("br"),n=_(`
                        Set to `),b=s("code"),b.textContent="null",p=_(" to delete already uploaded file(s).")},m(c,u){o(c,e,u),o(c,t,u),o(c,n,u),o(c,b,u),o(c,p,u)},p:te,d(c){c&&(d(e),d(t),d(n),d(b),d(p))}}}function jt(r){let e,t;return{c(){e=s("code"),e.textContent='{"lon":x,"lat":y}',t=_(" object.")},m(n,b){o(n,e,b),o(n,t,b)},p:te,d(n){n&&(d(e),d(t))}}}function Jt(r){let e;return{c(){e=_("URL address.")},m(t,n){o(t,e,n)},p:te,d(t){t&&d(e)}}}function Ut(r){let e;return{c(){e=_("Email address.")},m(t,n){o(t,e,n)},p:te,d(t){t&&d(e)}}}function Vt(r){let e;return{c(){e=_("JSON array or object.")},m(t,n){o(t,e,n)},p:te,d(t){t&&d(e)}}}function xt(r){let e;return{c(){e=_("Number value.")},m(t,n){o(t,e,n)},p:te,d(t){t&&d(e)}}}function Qt(r){let e;return{c(){e=_("Plain text value.")},m(t,n){o(t,e,n)},p:te,d(t){t&&d(e)}}}function gt(r,e){let t,n,b,p,c,u=e[15].name+"",m,O,T,H,P=ye.getFieldValueType(e[15])+"",$,M,q,L;function J(g,w){return g[15].required?At:Nt}let j=J(e),S=j(e);function R(g,w){if(g[15].type==="text")return Qt;if(g[15].type==="number")return xt;if(g[15].type==="json")return Vt;if(g[15].type==="email")return Ut;if(g[15].type==="url")return Jt;if(g[15].type==="geoPoint")return jt;if(g[15].type==="file")return It;if(g[15].type==="relation")return Et}let D=R(e),v=D&&D(e);return{key:r,first:null,c(){t=s("tr"),n=s("td"),b=s("div"),S.c(),p=f(),c=s("span"),m=_(u),O=f(),T=s("td"),H=s("span"),$=_(P),M=f(),q=s("td"),v&&v.c(),L=f(),k(b,"class","inline-flex"),k(H,"class","label"),this.first=t},m(g,w){o(g,t,w),a(t,n),a(n,b),S.m(b,null),a(b,p),a(b,c),a(c,m),a(t,O),a(t,T),a(T,H),a(H,$),a(t,M),a(t,q),v&&v.m(q,null),a(t,L)},p(g,w){e=g,j!==(j=J(e))&&(S.d(1),S=j(e),S&&(S.c(),S.m(b,p))),w&32&&u!==(u=e[15].name+"")&&ee(m,u),w&32&&P!==(P=ye.getFieldValueType(e[15])+"")&&ee($,P),D===(D=R(e))&&v?v.p(e,w):(v&&v.d(1),v=D&&D(e),v&&(v.c(),v.m(q,null)))},d(g){g&&d(t),S.d(),v&&v.d()}}}function wt(r,e){let t,n=e[10].code+"",b,p,c,u;function m(){return e[9](e[10])}return{key:r,first:null,c(){t=s("button"),b=_(n),p=f(),k(t,"class","tab-item"),Te(t,"active",e[2]===e[10].code),this.first=t},m(O,T){o(O,t,T),a(t,b),a(t,p),c||(u=Lt(t,"click",m),c=!0)},p(O,T){e=O,T&8&&n!==(n=e[10].code+"")&&ee(b,n),T&12&&Te(t,"active",e[2]===e[10].code)},d(O){O&&d(t),c=!1,u()}}}function Ct(r,e){let t,n,b,p;return n=new Tt({props:{content:e[10].body}}),{key:r,first:null,c(){t=s("div"),Ce(n.$$.fragment),b=f(),k(t,"class","tab-item"),Te(t,"active",e[2]===e[10].code),this.first=t},m(c,u){o(c,t,u),we(n,t,null),a(t,b),p=!0},p(c,u){e=c;const m={};u&8&&(m.content=e[10].body),n.$set(m),(!p||u&12)&&Te(t,"active",e[2]===e[10].code)},i(c){p||(he(n.$$.fragment,c),p=!0)},o(c){_e(n.$$.fragment,c),p=!1},d(c){c&&d(t),ge(n)}}}function zt(r){var ct,ut;let e,t,n=r[0].name+"",b,p,c,u,m,O,T,H=r[0].name+"",P,$,M,q,L,J,j,S,R,D,v,g,w,x,Q,B,h,C,le,z=r[0].name+"",ne,Ue,$e,Ve,Me,de,Oe,oe,Se,re,qe,K,Re,xe,W,He,U=[],Qe=new Map,Pe,ce,Le,Y,De,ze,ue,G,Fe,Ke,Be,We,N,Ye,ae,Ge,Xe,Ze,Ne,et,Ae,tt,Ee,lt,nt,ie,Ie,pe,je,X,fe,V=[],at=new Map,it,be,A=[],st=new Map,Z,E=r[1]&&yt();R=new St({props:{js:`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${r[4]}');

...

// example update data
const data = ${JSON.stringify(r[7](r[0]),null,4)};

const record = await pb.collection('${(ct=r[0])==null?void 0:ct.name}').update('RECORD_ID', data);
    `,dart:`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${r[4]}');

...

// example update body
final body = <String, dynamic>${JSON.stringify(r[7](r[0]),null,2)};

final record = await pb.collection('${(ut=r[0])==null?void 0:ut.name}').update('RECORD_ID', body: body);
    `}});let I=r[6]&&kt(),F=r[1]&&vt(r),ke=se(r[5]);const dt=l=>l[15].name;for(let l=0;l<ke.length;l+=1){let i=ht(r,ke,l),y=dt(i);Qe.set(y,U[l]=gt(y,i))}ae=new Tt({props:{content:"?expand=relField1,relField2.subRelField21"}}),ie=new Dt({});let ve=se(r[3]);const ot=l=>l[10].code;for(let l=0;l<ve.length;l+=1){let i=_t(r,ve,l),y=ot(i);at.set(y,V[l]=wt(y,i))}let me=se(r[3]);const rt=l=>l[10].code;for(let l=0;l<me.length;l+=1){let i=mt(r,me,l),y=rt(i);st.set(y,A[l]=Ct(y,i))}return{c(){e=s("h3"),t=_("Update ("),b=_(n),p=_(")"),c=f(),u=s("div"),m=s("p"),O=_("Update a single "),T=s("strong"),P=_(H),$=_(" record."),M=f(),q=s("p"),q.innerHTML=`Body parameters could be sent as <code>application/json</code> or
        <code>multipart/form-data</code>.`,L=f(),J=s("p"),J.innerHTML=`File upload is supported only via <code>multipart/form-data</code>.
        <br/>
        For more info and examples you could check the detailed
        <a href="https://pocketbase.io/docs/files-handling" target="_blank" rel="noopener noreferrer">Files upload and handling docs
        </a>.`,j=f(),E&&E.c(),S=f(),Ce(R.$$.fragment),D=f(),v=s("h6"),v.textContent="API details",g=f(),w=s("div"),x=s("strong"),x.textContent="PATCH",Q=f(),B=s("div"),h=s("p"),C=_("/api/collections/"),le=s("strong"),ne=_(z),Ue=_("/records/"),$e=s("strong"),$e.textContent=":id",Ve=f(),I&&I.c(),Me=f(),de=s("div"),de.textContent="Path parameters",Oe=f(),oe=s("table"),oe.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr></thead> <tbody><tr><td>id</td> <td><span class="label">String</span></td> <td>ID of the record to update.</td></tr></tbody>',Se=f(),re=s("div"),re.textContent="Body Parameters",qe=f(),K=s("table"),Re=s("thead"),Re.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr>',xe=f(),W=s("tbody"),F&&F.c(),He=f();for(let l=0;l<U.length;l+=1)U[l].c();Pe=f(),ce=s("div"),ce.textContent="Query parameters",Le=f(),Y=s("table"),De=s("thead"),De.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr>',ze=f(),ue=s("tbody"),G=s("tr"),Fe=s("td"),Fe.textContent="expand",Ke=f(),Be=s("td"),Be.innerHTML='<span class="label">String</span>',We=f(),N=s("td"),Ye=_(`Auto expand relations when returning the updated record. Ex.:
                `),Ce(ae.$$.fragment),Ge=_(`
                Supports up to 6-levels depth nested relations expansion. `),Xe=s("br"),Ze=_(`
                The expanded relations will be appended to the record under the
                `),Ne=s("code"),Ne.textContent="expand",et=_(" property (eg. "),Ae=s("code"),Ae.textContent='"expand": {"relField1": {...}, ...}',tt=_(`). Only
                the relations that the user has permissions to `),Ee=s("strong"),Ee.textContent="view",lt=_(" will be expanded."),nt=f(),Ce(ie.$$.fragment),Ie=f(),pe=s("div"),pe.textContent="Responses",je=f(),X=s("div"),fe=s("div");for(let l=0;l<V.length;l+=1)V[l].c();it=f(),be=s("div");for(let l=0;l<A.length;l+=1)A[l].c();k(e,"class","m-b-sm"),k(u,"class","content txt-lg m-b-sm"),k(v,"class","m-b-xs"),k(x,"class","label label-primary"),k(B,"class","content"),k(w,"class","alert alert-warning"),k(de,"class","section-title"),k(oe,"class","table-compact table-border m-b-base"),k(re,"class","section-title"),k(K,"class","table-compact table-border m-b-base"),k(ce,"class","section-title"),k(Y,"class","table-compact table-border m-b-lg"),k(pe,"class","section-title"),k(fe,"class","tabs-header compact combined left"),k(be,"class","tabs-content"),k(X,"class","tabs")},m(l,i){o(l,e,i),a(e,t),a(e,b),a(e,p),o(l,c,i),o(l,u,i),a(u,m),a(m,O),a(m,T),a(T,P),a(m,$),a(u,M),a(u,q),a(u,L),a(u,J),a(u,j),E&&E.m(u,null),o(l,S,i),we(R,l,i),o(l,D,i),o(l,v,i),o(l,g,i),o(l,w,i),a(w,x),a(w,Q),a(w,B),a(B,h),a(h,C),a(h,le),a(le,ne),a(h,Ue),a(h,$e),a(w,Ve),I&&I.m(w,null),o(l,Me,i),o(l,de,i),o(l,Oe,i),o(l,oe,i),o(l,Se,i),o(l,re,i),o(l,qe,i),o(l,K,i),a(K,Re),a(K,xe),a(K,W),F&&F.m(W,null),a(W,He);for(let y=0;y<U.length;y+=1)U[y]&&U[y].m(W,null);o(l,Pe,i),o(l,ce,i),o(l,Le,i),o(l,Y,i),a(Y,De),a(Y,ze),a(Y,ue),a(ue,G),a(G,Fe),a(G,Ke),a(G,Be),a(G,We),a(G,N),a(N,Ye),we(ae,N,null),a(N,Ge),a(N,Xe),a(N,Ze),a(N,Ne),a(N,et),a(N,Ae),a(N,tt),a(N,Ee),a(N,lt),a(ue,nt),we(ie,ue,null),o(l,Ie,i),o(l,pe,i),o(l,je,i),o(l,X,i),a(X,fe);for(let y=0;y<V.length;y+=1)V[y]&&V[y].m(fe,null);a(X,it),a(X,be);for(let y=0;y<A.length;y+=1)A[y]&&A[y].m(be,null);Z=!0},p(l,[i]){var pt,ft;(!Z||i&1)&&n!==(n=l[0].name+"")&&ee(b,n),(!Z||i&1)&&H!==(H=l[0].name+"")&&ee(P,H),l[1]?E||(E=yt(),E.c(),E.m(u,null)):E&&(E.d(1),E=null);const y={};i&17&&(y.js=`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${l[4]}');

...

// example update data
const data = ${JSON.stringify(l[7](l[0]),null,4)};

const record = await pb.collection('${(pt=l[0])==null?void 0:pt.name}').update('RECORD_ID', data);
    `),i&17&&(y.dart=`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${l[4]}');

...

// example update body
final body = <String, dynamic>${JSON.stringify(l[7](l[0]),null,2)};

final record = await pb.collection('${(ft=l[0])==null?void 0:ft.name}').update('RECORD_ID', body: body);
    `),R.$set(y),(!Z||i&1)&&z!==(z=l[0].name+"")&&ee(ne,z),l[6]?I||(I=kt(),I.c(),I.m(w,null)):I&&(I.d(1),I=null),l[1]?F?F.p(l,i):(F=vt(l),F.c(),F.m(W,He)):F&&(F.d(1),F=null),i&32&&(ke=se(l[5]),U=Je(U,i,dt,1,l,ke,Qe,W,bt,gt,null,ht)),i&12&&(ve=se(l[3]),V=Je(V,i,ot,1,l,ve,at,fe,bt,wt,null,_t)),i&12&&(me=se(l[3]),qt(),A=Je(A,i,rt,1,l,me,st,be,Rt,Ct,null,mt),Ht())},i(l){if(!Z){he(R.$$.fragment,l),he(ae.$$.fragment,l),he(ie.$$.fragment,l);for(let i=0;i<me.length;i+=1)he(A[i]);Z=!0}},o(l){_e(R.$$.fragment,l),_e(ae.$$.fragment,l),_e(ie.$$.fragment,l);for(let i=0;i<A.length;i+=1)_e(A[i]);Z=!1},d(l){l&&(d(e),d(c),d(u),d(S),d(D),d(v),d(g),d(w),d(Me),d(de),d(Oe),d(oe),d(Se),d(re),d(qe),d(K),d(Pe),d(ce),d(Le),d(Y),d(Ie),d(pe),d(je),d(X)),E&&E.d(),ge(R,l),I&&I.d(),F&&F.d();for(let i=0;i<U.length;i+=1)U[i].d();ge(ae),ge(ie);for(let i=0;i<V.length;i+=1)V[i].d();for(let i=0;i<A.length;i+=1)A[i].d()}}}const Kt=r=>r.name=="emailVisibility";function Wt(r,e,t){let n,b,p,c,u,{collection:m}=e,O=200,T=[];function H($){let M=ye.dummyCollectionSchemaData($,!0);return n&&(M.oldPassword="12345678",M.password="87654321",M.passwordConfirm="87654321",delete M.verified,delete M.email),M}const P=$=>t(2,O=$.code);return r.$$set=$=>{"collection"in $&&t(0,m=$.collection)},r.$$.update=()=>{var $,M,q;r.$$.dirty&1&&t(1,n=(m==null?void 0:m.type)==="auth"),r.$$.dirty&1&&t(6,b=(m==null?void 0:m.updateRule)===null),r.$$.dirty&2&&t(8,p=n?["id","password","verified","email","emailVisibility"]:["id"]),r.$$.dirty&257&&t(5,c=(($=m==null?void 0:m.fields)==null?void 0:$.filter(L=>!L.hidden&&L.type!="autodate"&&!p.includes(L.name)))||[]),r.$$.dirty&1&&t(3,T=[{code:200,body:JSON.stringify(ye.dummyCollectionRecord(m),null,2)},{code:400,body:`
                {
                  "status": 400,
                  "message": "Failed to update record.",
                  "data": {
                    "${(q=(M=m==null?void 0:m.fields)==null?void 0:M[0])==null?void 0:q.name}": {
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
            `},{code:404,body:`
                {
                  "status": 404,
                  "message": "The requested resource wasn't found.",
                  "data": {}
                }
            `}])},t(4,u=ye.getApiExampleUrl(Pt.baseURL)),[m,n,O,T,u,c,b,H,p,P]}class Xt extends $t{constructor(e){super(),Mt(this,e,Wt,zt,Ot,{collection:0})}}export{Xt as default};
