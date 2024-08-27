import{S as qt,i as Ot,s as Mt,C as Q,O as ne,N as Tt,e as s,v as _,b as f,c as _e,f as v,g as r,h as n,m as he,w as x,P as Be,Q as ht,k as Ht,R as Lt,n as Pt,t as fe,a as ue,o as d,d as ke,A as At,q as ye,r as Ft,x as ae}from"./index-D0DO79Dq.js";import{S as Rt}from"./SdkTabs-DC6EUYpr.js";import{F as Bt}from"./FieldsQueryParam-BwleQAus.js";function kt(o,e,t){const a=o.slice();return a[8]=e[t],a}function yt(o,e,t){const a=o.slice();return a[8]=e[t],a}function vt(o,e,t){const a=o.slice();return a[13]=e[t],a}function gt(o){let e;return{c(){e=s("p"),e.innerHTML="Requires admin <code>Authorization:TOKEN</code> header",v(e,"class","txt-hint txt-sm txt-right")},m(t,a){r(t,e,a)},d(t){t&&d(e)}}}function wt(o){let e,t,a,u,m,c,p,y,S,T,w,H,D,E,P,I,j,B,C,N,q,g,b;function O(h,$){var ee,K;return(K=(ee=h[0])==null?void 0:ee.options)!=null&&K.requireEmail?Dt:jt}let z=O(o),A=z(o);return{c(){e=s("tr"),e.innerHTML='<td colspan="3" class="txt-hint">Auth fields</td>',t=f(),a=s("tr"),a.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>username</span></div></td> <td><span class="label">String</span></td> <td>The username of the auth record.
                    <br/>
                    If not set, it will be auto generated.</td>`,u=f(),m=s("tr"),c=s("td"),p=s("div"),A.c(),y=f(),S=s("span"),S.textContent="email",T=f(),w=s("td"),w.innerHTML='<span class="label">String</span>',H=f(),D=s("td"),D.textContent="Auth record email address.",E=f(),P=s("tr"),P.innerHTML='<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>emailVisibility</span></div></td> <td><span class="label">Boolean</span></td> <td>Whether to show/hide the auth record email when fetching the record data.</td>',I=f(),j=s("tr"),j.innerHTML='<td><div class="inline-flex"><span class="label label-success">Required</span> <span>password</span></div></td> <td><span class="label">String</span></td> <td>Auth record password.</td>',B=f(),C=s("tr"),C.innerHTML='<td><div class="inline-flex"><span class="label label-success">Required</span> <span>passwordConfirm</span></div></td> <td><span class="label">String</span></td> <td>Auth record password confirmation.</td>',N=f(),q=s("tr"),q.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>verified</span></div></td> <td><span class="label">Boolean</span></td> <td>Indicates whether the auth record is verified or not.
                    <br/>
                    This field can be set only by admins or auth records with &quot;Manage&quot; access.</td>`,g=f(),b=s("tr"),b.innerHTML='<td colspan="3" class="txt-hint">Schema fields</td>',v(p,"class","inline-flex")},m(h,$){r(h,e,$),r(h,t,$),r(h,a,$),r(h,u,$),r(h,m,$),n(m,c),n(c,p),A.m(p,null),n(p,y),n(p,S),n(m,T),n(m,w),n(m,H),n(m,D),r(h,E,$),r(h,P,$),r(h,I,$),r(h,j,$),r(h,B,$),r(h,C,$),r(h,N,$),r(h,q,$),r(h,g,$),r(h,b,$)},p(h,$){z!==(z=O(h))&&(A.d(1),A=z(h),A&&(A.c(),A.m(p,y)))},d(h){h&&(d(e),d(t),d(a),d(u),d(m),d(E),d(P),d(I),d(j),d(B),d(C),d(N),d(q),d(g),d(b)),A.d()}}}function jt(o){let e;return{c(){e=s("span"),e.textContent="Optional",v(e,"class","label label-warning")},m(t,a){r(t,e,a)},d(t){t&&d(e)}}}function Dt(o){let e;return{c(){e=s("span"),e.textContent="Required",v(e,"class","label label-success")},m(t,a){r(t,e,a)},d(t){t&&d(e)}}}function Nt(o){let e;return{c(){e=s("span"),e.textContent="Optional",v(e,"class","label label-warning")},m(t,a){r(t,e,a)},d(t){t&&d(e)}}}function Vt(o){let e;return{c(){e=s("span"),e.textContent="Required",v(e,"class","label label-success")},m(t,a){r(t,e,a)},d(t){t&&d(e)}}}function Jt(o){var m;let e,t=((m=o[13].options)==null?void 0:m.maxSelect)===1?"id":"ids",a,u;return{c(){e=_("Relation record "),a=_(t),u=_(".")},m(c,p){r(c,e,p),r(c,a,p),r(c,u,p)},p(c,p){var y;p&1&&t!==(t=((y=c[13].options)==null?void 0:y.maxSelect)===1?"id":"ids")&&x(a,t)},d(c){c&&(d(e),d(a),d(u))}}}function Et(o){let e,t,a,u,m;return{c(){e=_("File object."),t=s("br"),a=_(`
                        Set to `),u=s("code"),u.textContent="null",m=_(" to delete already uploaded file(s).")},m(c,p){r(c,e,p),r(c,t,p),r(c,a,p),r(c,u,p),r(c,m,p)},p:ae,d(c){c&&(d(e),d(t),d(a),d(u),d(m))}}}function It(o){let e;return{c(){e=_("URL address.")},m(t,a){r(t,e,a)},p:ae,d(t){t&&d(e)}}}function Ut(o){let e;return{c(){e=_("Email address.")},m(t,a){r(t,e,a)},p:ae,d(t){t&&d(e)}}}function Qt(o){let e;return{c(){e=_("JSON array or object.")},m(t,a){r(t,e,a)},p:ae,d(t){t&&d(e)}}}function zt(o){let e;return{c(){e=_("Number value.")},m(t,a){r(t,e,a)},p:ae,d(t){t&&d(e)}}}function Kt(o){let e;return{c(){e=_("Plain text value.")},m(t,a){r(t,e,a)},p:ae,d(t){t&&d(e)}}}function Ct(o,e){let t,a,u,m,c,p=e[13].name+"",y,S,T,w,H=Q.getFieldValueType(e[13])+"",D,E,P,I;function j(b,O){return b[13].required?Vt:Nt}let B=j(e),C=B(e);function N(b,O){if(b[13].type==="text")return Kt;if(b[13].type==="number")return zt;if(b[13].type==="json")return Qt;if(b[13].type==="email")return Ut;if(b[13].type==="url")return It;if(b[13].type==="file")return Et;if(b[13].type==="relation")return Jt}let q=N(e),g=q&&q(e);return{key:o,first:null,c(){t=s("tr"),a=s("td"),u=s("div"),C.c(),m=f(),c=s("span"),y=_(p),S=f(),T=s("td"),w=s("span"),D=_(H),E=f(),P=s("td"),g&&g.c(),I=f(),v(u,"class","inline-flex"),v(w,"class","label"),this.first=t},m(b,O){r(b,t,O),n(t,a),n(a,u),C.m(u,null),n(u,m),n(u,c),n(c,y),n(t,S),n(t,T),n(T,w),n(w,D),n(t,E),n(t,P),g&&g.m(P,null),n(t,I)},p(b,O){e=b,B!==(B=j(e))&&(C.d(1),C=B(e),C&&(C.c(),C.m(u,m))),O&1&&p!==(p=e[13].name+"")&&x(y,p),O&1&&H!==(H=Q.getFieldValueType(e[13])+"")&&x(D,H),q===(q=N(e))&&g?g.p(e,O):(g&&g.d(1),g=q&&q(e),g&&(g.c(),g.m(P,null)))},d(b){b&&d(t),C.d(),g&&g.d()}}}function $t(o,e){let t,a=e[8].code+"",u,m,c,p;function y(){return e[7](e[8])}return{key:o,first:null,c(){t=s("button"),u=_(a),m=f(),v(t,"class","tab-item"),ye(t,"active",e[2]===e[8].code),this.first=t},m(S,T){r(S,t,T),n(t,u),n(t,m),c||(p=Ft(t,"click",y),c=!0)},p(S,T){e=S,T&8&&a!==(a=e[8].code+"")&&x(u,a),T&12&&ye(t,"active",e[2]===e[8].code)},d(S){S&&d(t),c=!1,p()}}}function St(o,e){let t,a,u,m;return a=new Tt({props:{content:e[8].body}}),{key:o,first:null,c(){t=s("div"),_e(a.$$.fragment),u=f(),v(t,"class","tab-item"),ye(t,"active",e[2]===e[8].code),this.first=t},m(c,p){r(c,t,p),he(a,t,null),n(t,u),m=!0},p(c,p){e=c;const y={};p&8&&(y.content=e[8].body),a.$set(y),(!m||p&12)&&ye(t,"active",e[2]===e[8].code)},i(c){m||(fe(a.$$.fragment,c),m=!0)},o(c){ue(a.$$.fragment,c),m=!1},d(c){c&&d(t),ke(a)}}}function Wt(o){var ot,rt,dt,ct,pt;let e,t,a=o[0].name+"",u,m,c,p,y,S,T,w=o[0].name+"",H,D,E,P,I,j,B,C,N,q,g,b,O,z,A,h,$,ee,K=o[0].name+"",ve,je,De,ge,ie,we,W,Ce,Ne,U,$e,Ve,Se,V=[],Je=new Map,Te,se,qe,Y,Oe,Ee,oe,G,Me,Ie,He,Ue,M,Qe,te,ze,Ke,We,Le,Ye,Pe,Ge,Xe,Ze,Ae,xe,et,le,Fe,re,Re,X,de,J=[],tt=new Map,lt,ce,F=[],nt=new Map,Z;C=new Rt({props:{js:`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${o[5]}');

...

// example create data
const data = ${JSON.stringify(Object.assign({},o[4],Q.dummyCollectionSchemaData(o[0])),null,4)};

const record = await pb.collection('${(ot=o[0])==null?void 0:ot.name}').create(data);
`+(o[1]?`
// (optional) send an email verification request
await pb.collection('${(rt=o[0])==null?void 0:rt.name}').requestVerification('test@example.com');
`:""),dart:`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${o[5]}');

...

// example create body
final body = <String, dynamic>${JSON.stringify(Object.assign({},o[4],Q.dummyCollectionSchemaData(o[0])),null,2)};

final record = await pb.collection('${(dt=o[0])==null?void 0:dt.name}').create(body: body);
`+(o[1]?`
// (optional) send an email verification request
await pb.collection('${(ct=o[0])==null?void 0:ct.name}').requestVerification('test@example.com');
`:"")}});let R=o[6]&&gt(),L=o[1]&&wt(o),me=ne((pt=o[0])==null?void 0:pt.schema);const at=l=>l[13].name;for(let l=0;l<me.length;l+=1){let i=vt(o,me,l),k=at(i);Je.set(k,V[l]=Ct(k,i))}te=new Tt({props:{content:"?expand=relField1,relField2.subRelField"}}),le=new Bt({});let be=ne(o[3]);const it=l=>l[8].code;for(let l=0;l<be.length;l+=1){let i=yt(o,be,l),k=it(i);tt.set(k,J[l]=$t(k,i))}let pe=ne(o[3]);const st=l=>l[8].code;for(let l=0;l<pe.length;l+=1){let i=kt(o,pe,l),k=st(i);nt.set(k,F[l]=St(k,i))}return{c(){e=s("h3"),t=_("Create ("),u=_(a),m=_(")"),c=f(),p=s("div"),y=s("p"),S=_("Create a new "),T=s("strong"),H=_(w),D=_(" record."),E=f(),P=s("p"),P.innerHTML=`Body parameters could be sent as <code>application/json</code> or
        <code>multipart/form-data</code>.`,I=f(),j=s("p"),j.innerHTML=`File upload is supported only via <code>multipart/form-data</code>.
        <br/>
        For more info and examples you could check the detailed
        <a href="https://pocketbase.io/docs/files-handling/" target="_blank" rel="noopener noreferrer">Files upload and handling docs
        </a>.`,B=f(),_e(C.$$.fragment),N=f(),q=s("h6"),q.textContent="API details",g=f(),b=s("div"),O=s("strong"),O.textContent="POST",z=f(),A=s("div"),h=s("p"),$=_("/api/collections/"),ee=s("strong"),ve=_(K),je=_("/records"),De=f(),R&&R.c(),ge=f(),ie=s("div"),ie.textContent="Body Parameters",we=f(),W=s("table"),Ce=s("thead"),Ce.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr>',Ne=f(),U=s("tbody"),$e=s("tr"),$e.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>id</span></div></td> <td><span class="label">String</span></td> <td><strong>15 characters string</strong> to store as record ID.
                <br/>
                If not set, it will be auto generated.</td>`,Ve=f(),L&&L.c(),Se=f();for(let l=0;l<V.length;l+=1)V[l].c();Te=f(),se=s("div"),se.textContent="Query parameters",qe=f(),Y=s("table"),Oe=s("thead"),Oe.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr>',Ee=f(),oe=s("tbody"),G=s("tr"),Me=s("td"),Me.textContent="expand",Ie=f(),He=s("td"),He.innerHTML='<span class="label">String</span>',Ue=f(),M=s("td"),Qe=_(`Auto expand relations when returning the created record. Ex.:
                `),_e(te.$$.fragment),ze=_(`
                Supports up to 6-levels depth nested relations expansion. `),Ke=s("br"),We=_(`
                The expanded relations will be appended to the record under the
                `),Le=s("code"),Le.textContent="expand",Ye=_(" property (eg. "),Pe=s("code"),Pe.textContent='"expand": {"relField1": {...}, ...}',Ge=_(`).
                `),Xe=s("br"),Ze=_(`
                Only the relations to which the request user has permissions to `),Ae=s("strong"),Ae.textContent="view",xe=_(" will be expanded."),et=f(),_e(le.$$.fragment),Fe=f(),re=s("div"),re.textContent="Responses",Re=f(),X=s("div"),de=s("div");for(let l=0;l<J.length;l+=1)J[l].c();lt=f(),ce=s("div");for(let l=0;l<F.length;l+=1)F[l].c();v(e,"class","m-b-sm"),v(p,"class","content txt-lg m-b-sm"),v(q,"class","m-b-xs"),v(O,"class","label label-primary"),v(A,"class","content"),v(b,"class","alert alert-success"),v(ie,"class","section-title"),v(W,"class","table-compact table-border m-b-base"),v(se,"class","section-title"),v(Y,"class","table-compact table-border m-b-base"),v(re,"class","section-title"),v(de,"class","tabs-header compact combined left"),v(ce,"class","tabs-content"),v(X,"class","tabs")},m(l,i){r(l,e,i),n(e,t),n(e,u),n(e,m),r(l,c,i),r(l,p,i),n(p,y),n(y,S),n(y,T),n(T,H),n(y,D),n(p,E),n(p,P),n(p,I),n(p,j),r(l,B,i),he(C,l,i),r(l,N,i),r(l,q,i),r(l,g,i),r(l,b,i),n(b,O),n(b,z),n(b,A),n(A,h),n(h,$),n(h,ee),n(ee,ve),n(h,je),n(b,De),R&&R.m(b,null),r(l,ge,i),r(l,ie,i),r(l,we,i),r(l,W,i),n(W,Ce),n(W,Ne),n(W,U),n(U,$e),n(U,Ve),L&&L.m(U,null),n(U,Se);for(let k=0;k<V.length;k+=1)V[k]&&V[k].m(U,null);r(l,Te,i),r(l,se,i),r(l,qe,i),r(l,Y,i),n(Y,Oe),n(Y,Ee),n(Y,oe),n(oe,G),n(G,Me),n(G,Ie),n(G,He),n(G,Ue),n(G,M),n(M,Qe),he(te,M,null),n(M,ze),n(M,Ke),n(M,We),n(M,Le),n(M,Ye),n(M,Pe),n(M,Ge),n(M,Xe),n(M,Ze),n(M,Ae),n(M,xe),n(oe,et),he(le,oe,null),r(l,Fe,i),r(l,re,i),r(l,Re,i),r(l,X,i),n(X,de);for(let k=0;k<J.length;k+=1)J[k]&&J[k].m(de,null);n(X,lt),n(X,ce);for(let k=0;k<F.length;k+=1)F[k]&&F[k].m(ce,null);Z=!0},p(l,[i]){var ft,ut,mt,bt,_t;(!Z||i&1)&&a!==(a=l[0].name+"")&&x(u,a),(!Z||i&1)&&w!==(w=l[0].name+"")&&x(H,w);const k={};i&51&&(k.js=`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${l[5]}');

...

// example create data
const data = ${JSON.stringify(Object.assign({},l[4],Q.dummyCollectionSchemaData(l[0])),null,4)};

const record = await pb.collection('${(ft=l[0])==null?void 0:ft.name}').create(data);
`+(l[1]?`
// (optional) send an email verification request
await pb.collection('${(ut=l[0])==null?void 0:ut.name}').requestVerification('test@example.com');
`:"")),i&51&&(k.dart=`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${l[5]}');

...

// example create body
final body = <String, dynamic>${JSON.stringify(Object.assign({},l[4],Q.dummyCollectionSchemaData(l[0])),null,2)};

final record = await pb.collection('${(mt=l[0])==null?void 0:mt.name}').create(body: body);
`+(l[1]?`
// (optional) send an email verification request
await pb.collection('${(bt=l[0])==null?void 0:bt.name}').requestVerification('test@example.com');
`:"")),C.$set(k),(!Z||i&1)&&K!==(K=l[0].name+"")&&x(ve,K),l[6]?R||(R=gt(),R.c(),R.m(b,null)):R&&(R.d(1),R=null),l[1]?L?L.p(l,i):(L=wt(l),L.c(),L.m(U,Se)):L&&(L.d(1),L=null),i&1&&(me=ne((_t=l[0])==null?void 0:_t.schema),V=Be(V,i,at,1,l,me,Je,U,ht,Ct,null,vt)),i&12&&(be=ne(l[3]),J=Be(J,i,it,1,l,be,tt,de,ht,$t,null,yt)),i&12&&(pe=ne(l[3]),Ht(),F=Be(F,i,st,1,l,pe,nt,ce,Lt,St,null,kt),Pt())},i(l){if(!Z){fe(C.$$.fragment,l),fe(te.$$.fragment,l),fe(le.$$.fragment,l);for(let i=0;i<pe.length;i+=1)fe(F[i]);Z=!0}},o(l){ue(C.$$.fragment,l),ue(te.$$.fragment,l),ue(le.$$.fragment,l);for(let i=0;i<F.length;i+=1)ue(F[i]);Z=!1},d(l){l&&(d(e),d(c),d(p),d(B),d(N),d(q),d(g),d(b),d(ge),d(ie),d(we),d(W),d(Te),d(se),d(qe),d(Y),d(Fe),d(re),d(Re),d(X)),ke(C,l),R&&R.d(),L&&L.d();for(let i=0;i<V.length;i+=1)V[i].d();ke(te),ke(le);for(let i=0;i<J.length;i+=1)J[i].d();for(let i=0;i<F.length;i+=1)F[i].d()}}}function Yt(o,e,t){let a,u,m,{collection:c}=e,p=200,y=[],S={};const T=w=>t(2,p=w.code);return o.$$set=w=>{"collection"in w&&t(0,c=w.collection)},o.$$.update=()=>{var w,H;o.$$.dirty&1&&t(1,a=c.type==="auth"),o.$$.dirty&1&&t(6,u=(c==null?void 0:c.createRule)===null),o.$$.dirty&1&&t(3,y=[{code:200,body:JSON.stringify(Q.dummyCollectionRecord(c),null,2)},{code:400,body:`
                {
                  "code": 400,
                  "message": "Failed to create record.",
                  "data": {
                    "${(H=(w=c==null?void 0:c.schema)==null?void 0:w[0])==null?void 0:H.name}": {
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
            `}]),o.$$.dirty&2&&(a?t(4,S={username:"test_username",email:"test@example.com",emailVisibility:!0,password:"12345678",passwordConfirm:"12345678"}):t(4,S={}))},t(5,m=Q.getApiExampleUrl(At.baseUrl)),[c,a,p,y,S,m,u,T]}class xt extends qt{constructor(e){super(),Ot(this,e,Yt,Wt,Mt,{collection:0})}}export{xt as default};
