import{S as qt,i as Ot,s as Mt,C as Q,O as ne,N as Tt,e as i,w as _,b as u,c as _e,f as v,g as r,h as n,m as he,x,P as Be,Q as ht,k as Ht,R as Lt,n as Pt,t as ue,a as fe,o as d,d as ke,p as Ft,r as ye,u as Rt,y as ae}from"./index-78piLIP3.js";import{S as At}from"./SdkTabs-c6VuPJvR.js";import{F as Bt}from"./FieldsQueryParam-oYzijp1d.js";function kt(o,e,t){const a=o.slice();return a[8]=e[t],a}function yt(o,e,t){const a=o.slice();return a[8]=e[t],a}function vt(o,e,t){const a=o.slice();return a[13]=e[t],a}function gt(o){let e;return{c(){e=i("p"),e.innerHTML="Requires admin <code>Authorization:TOKEN</code> header",v(e,"class","txt-hint txt-sm txt-right")},m(t,a){r(t,e,a)},d(t){t&&d(e)}}}function wt(o){let e,t,a,f,m,c,p,y,S,T,w,H,D,E,P,I,j,B,$,N,q,g,b;function O(h,C){var ee,K;return(K=(ee=h[0])==null?void 0:ee.options)!=null&&K.requireEmail?Dt:jt}let z=O(o),F=z(o);return{c(){e=i("tr"),e.innerHTML='<td colspan="3" class="txt-hint">Auth fields</td>',t=u(),a=i("tr"),a.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>username</span></div></td> <td><span class="label">String</span></td> <td>The username of the auth record.
                    <br/>
                    If not set, it will be auto generated.</td>`,f=u(),m=i("tr"),c=i("td"),p=i("div"),F.c(),y=u(),S=i("span"),S.textContent="email",T=u(),w=i("td"),w.innerHTML='<span class="label">String</span>',H=u(),D=i("td"),D.textContent="Auth record email address.",E=u(),P=i("tr"),P.innerHTML='<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>emailVisibility</span></div></td> <td><span class="label">Boolean</span></td> <td>Whether to show/hide the auth record email when fetching the record data.</td>',I=u(),j=i("tr"),j.innerHTML='<td><div class="inline-flex"><span class="label label-success">Required</span> <span>password</span></div></td> <td><span class="label">String</span></td> <td>Auth record password.</td>',B=u(),$=i("tr"),$.innerHTML='<td><div class="inline-flex"><span class="label label-success">Required</span> <span>passwordConfirm</span></div></td> <td><span class="label">String</span></td> <td>Auth record password confirmation.</td>',N=u(),q=i("tr"),q.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>verified</span></div></td> <td><span class="label">Boolean</span></td> <td>Indicates whether the auth record is verified or not.
                    <br/>
                    This field can be set only by admins or auth records with &quot;Manage&quot; access.</td>`,g=u(),b=i("tr"),b.innerHTML='<td colspan="3" class="txt-hint">Schema fields</td>',v(p,"class","inline-flex")},m(h,C){r(h,e,C),r(h,t,C),r(h,a,C),r(h,f,C),r(h,m,C),n(m,c),n(c,p),F.m(p,null),n(p,y),n(p,S),n(m,T),n(m,w),n(m,H),n(m,D),r(h,E,C),r(h,P,C),r(h,I,C),r(h,j,C),r(h,B,C),r(h,$,C),r(h,N,C),r(h,q,C),r(h,g,C),r(h,b,C)},p(h,C){z!==(z=O(h))&&(F.d(1),F=z(h),F&&(F.c(),F.m(p,y)))},d(h){h&&(d(e),d(t),d(a),d(f),d(m),d(E),d(P),d(I),d(j),d(B),d($),d(N),d(q),d(g),d(b)),F.d()}}}function jt(o){let e;return{c(){e=i("span"),e.textContent="Optional",v(e,"class","label label-warning")},m(t,a){r(t,e,a)},d(t){t&&d(e)}}}function Dt(o){let e;return{c(){e=i("span"),e.textContent="Required",v(e,"class","label label-success")},m(t,a){r(t,e,a)},d(t){t&&d(e)}}}function Nt(o){let e;return{c(){e=i("span"),e.textContent="Optional",v(e,"class","label label-warning")},m(t,a){r(t,e,a)},d(t){t&&d(e)}}}function Vt(o){let e;return{c(){e=i("span"),e.textContent="Required",v(e,"class","label label-success")},m(t,a){r(t,e,a)},d(t){t&&d(e)}}}function Jt(o){var m;let e,t=((m=o[13].options)==null?void 0:m.maxSelect)===1?"id":"ids",a,f;return{c(){e=_("Relation record "),a=_(t),f=_(".")},m(c,p){r(c,e,p),r(c,a,p),r(c,f,p)},p(c,p){var y;p&1&&t!==(t=((y=c[13].options)==null?void 0:y.maxSelect)===1?"id":"ids")&&x(a,t)},d(c){c&&(d(e),d(a),d(f))}}}function Et(o){let e,t,a,f,m;return{c(){e=_("File object."),t=i("br"),a=_(`
                        Set to `),f=i("code"),f.textContent="null",m=_(" to delete already uploaded file(s).")},m(c,p){r(c,e,p),r(c,t,p),r(c,a,p),r(c,f,p),r(c,m,p)},p:ae,d(c){c&&(d(e),d(t),d(a),d(f),d(m))}}}function It(o){let e;return{c(){e=_("URL address.")},m(t,a){r(t,e,a)},p:ae,d(t){t&&d(e)}}}function Ut(o){let e;return{c(){e=_("Email address.")},m(t,a){r(t,e,a)},p:ae,d(t){t&&d(e)}}}function Qt(o){let e;return{c(){e=_("JSON array or object.")},m(t,a){r(t,e,a)},p:ae,d(t){t&&d(e)}}}function zt(o){let e;return{c(){e=_("Number value.")},m(t,a){r(t,e,a)},p:ae,d(t){t&&d(e)}}}function Kt(o){let e;return{c(){e=_("Plain text value.")},m(t,a){r(t,e,a)},p:ae,d(t){t&&d(e)}}}function $t(o,e){let t,a,f,m,c,p=e[13].name+"",y,S,T,w,H=Q.getFieldValueType(e[13])+"",D,E,P,I;function j(b,O){return b[13].required?Vt:Nt}let B=j(e),$=B(e);function N(b,O){if(b[13].type==="text")return Kt;if(b[13].type==="number")return zt;if(b[13].type==="json")return Qt;if(b[13].type==="email")return Ut;if(b[13].type==="url")return It;if(b[13].type==="file")return Et;if(b[13].type==="relation")return Jt}let q=N(e),g=q&&q(e);return{key:o,first:null,c(){t=i("tr"),a=i("td"),f=i("div"),$.c(),m=u(),c=i("span"),y=_(p),S=u(),T=i("td"),w=i("span"),D=_(H),E=u(),P=i("td"),g&&g.c(),I=u(),v(f,"class","inline-flex"),v(w,"class","label"),this.first=t},m(b,O){r(b,t,O),n(t,a),n(a,f),$.m(f,null),n(f,m),n(f,c),n(c,y),n(t,S),n(t,T),n(T,w),n(w,D),n(t,E),n(t,P),g&&g.m(P,null),n(t,I)},p(b,O){e=b,B!==(B=j(e))&&($.d(1),$=B(e),$&&($.c(),$.m(f,m))),O&1&&p!==(p=e[13].name+"")&&x(y,p),O&1&&H!==(H=Q.getFieldValueType(e[13])+"")&&x(D,H),q===(q=N(e))&&g?g.p(e,O):(g&&g.d(1),g=q&&q(e),g&&(g.c(),g.m(P,null)))},d(b){b&&d(t),$.d(),g&&g.d()}}}function Ct(o,e){let t,a=e[8].code+"",f,m,c,p;function y(){return e[7](e[8])}return{key:o,first:null,c(){t=i("button"),f=_(a),m=u(),v(t,"class","tab-item"),ye(t,"active",e[2]===e[8].code),this.first=t},m(S,T){r(S,t,T),n(t,f),n(t,m),c||(p=Rt(t,"click",y),c=!0)},p(S,T){e=S,T&8&&a!==(a=e[8].code+"")&&x(f,a),T&12&&ye(t,"active",e[2]===e[8].code)},d(S){S&&d(t),c=!1,p()}}}function St(o,e){let t,a,f,m;return a=new Tt({props:{content:e[8].body}}),{key:o,first:null,c(){t=i("div"),_e(a.$$.fragment),f=u(),v(t,"class","tab-item"),ye(t,"active",e[2]===e[8].code),this.first=t},m(c,p){r(c,t,p),he(a,t,null),n(t,f),m=!0},p(c,p){e=c;const y={};p&8&&(y.content=e[8].body),a.$set(y),(!m||p&12)&&ye(t,"active",e[2]===e[8].code)},i(c){m||(ue(a.$$.fragment,c),m=!0)},o(c){fe(a.$$.fragment,c),m=!1},d(c){c&&d(t),ke(a)}}}function Wt(o){var ot,rt,dt,ct,pt;let e,t,a=o[0].name+"",f,m,c,p,y,S,T,w=o[0].name+"",H,D,E,P,I,j,B,$,N,q,g,b,O,z,F,h,C,ee,K=o[0].name+"",ve,je,De,ge,se,we,W,$e,Ne,U,Ce,Ve,Se,V=[],Je=new Map,Te,ie,qe,Y,Oe,Ee,oe,G,Me,Ie,He,Ue,M,Qe,te,ze,Ke,We,Le,Ye,Pe,Ge,Xe,Ze,Fe,xe,et,le,Re,re,Ae,X,de,J=[],tt=new Map,lt,ce,R=[],nt=new Map,Z;$=new At({props:{js:`
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
`:"")}});let A=o[6]&&gt(),L=o[1]&&wt(o),me=ne((pt=o[0])==null?void 0:pt.schema);const at=l=>l[13].name;for(let l=0;l<me.length;l+=1){let s=vt(o,me,l),k=at(s);Je.set(k,V[l]=$t(k,s))}te=new Tt({props:{content:"?expand=relField1,relField2.subRelField"}}),le=new Bt({});let be=ne(o[3]);const st=l=>l[8].code;for(let l=0;l<be.length;l+=1){let s=yt(o,be,l),k=st(s);tt.set(k,J[l]=Ct(k,s))}let pe=ne(o[3]);const it=l=>l[8].code;for(let l=0;l<pe.length;l+=1){let s=kt(o,pe,l),k=it(s);nt.set(k,R[l]=St(k,s))}return{c(){e=i("h3"),t=_("Create ("),f=_(a),m=_(")"),c=u(),p=i("div"),y=i("p"),S=_("Create a new "),T=i("strong"),H=_(w),D=_(" record."),E=u(),P=i("p"),P.innerHTML=`Body parameters could be sent as <code>application/json</code> or
        <code>multipart/form-data</code>.`,I=u(),j=i("p"),j.innerHTML=`File upload is supported only via <code>multipart/form-data</code>.
        <br/>
        For more info and examples you could check the detailed
        <a href="https://pocketbase.io/docs/files-handling/" target="_blank" rel="noopener noreferrer">Files upload and handling docs
        </a>.`,B=u(),_e($.$$.fragment),N=u(),q=i("h6"),q.textContent="API details",g=u(),b=i("div"),O=i("strong"),O.textContent="POST",z=u(),F=i("div"),h=i("p"),C=_("/api/collections/"),ee=i("strong"),ve=_(K),je=_("/records"),De=u(),A&&A.c(),ge=u(),se=i("div"),se.textContent="Body Parameters",we=u(),W=i("table"),$e=i("thead"),$e.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr>',Ne=u(),U=i("tbody"),Ce=i("tr"),Ce.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>id</span></div></td> <td><span class="label">String</span></td> <td><strong>15 characters string</strong> to store as record ID.
                <br/>
                If not set, it will be auto generated.</td>`,Ve=u(),L&&L.c(),Se=u();for(let l=0;l<V.length;l+=1)V[l].c();Te=u(),ie=i("div"),ie.textContent="Query parameters",qe=u(),Y=i("table"),Oe=i("thead"),Oe.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr>',Ee=u(),oe=i("tbody"),G=i("tr"),Me=i("td"),Me.textContent="expand",Ie=u(),He=i("td"),He.innerHTML='<span class="label">String</span>',Ue=u(),M=i("td"),Qe=_(`Auto expand relations when returning the created record. Ex.:
                `),_e(te.$$.fragment),ze=_(`
                Supports up to 6-levels depth nested relations expansion. `),Ke=i("br"),We=_(`
                The expanded relations will be appended to the record under the
                `),Le=i("code"),Le.textContent="expand",Ye=_(" property (eg. "),Pe=i("code"),Pe.textContent='"expand": {"relField1": {...}, ...}',Ge=_(`).
                `),Xe=i("br"),Ze=_(`
                Only the relations to which the request user has permissions to `),Fe=i("strong"),Fe.textContent="view",xe=_(" will be expanded."),et=u(),_e(le.$$.fragment),Re=u(),re=i("div"),re.textContent="Responses",Ae=u(),X=i("div"),de=i("div");for(let l=0;l<J.length;l+=1)J[l].c();lt=u(),ce=i("div");for(let l=0;l<R.length;l+=1)R[l].c();v(e,"class","m-b-sm"),v(p,"class","content txt-lg m-b-sm"),v(q,"class","m-b-xs"),v(O,"class","label label-primary"),v(F,"class","content"),v(b,"class","alert alert-success"),v(se,"class","section-title"),v(W,"class","table-compact table-border m-b-base"),v(ie,"class","section-title"),v(Y,"class","table-compact table-border m-b-base"),v(re,"class","section-title"),v(de,"class","tabs-header compact combined left"),v(ce,"class","tabs-content"),v(X,"class","tabs")},m(l,s){r(l,e,s),n(e,t),n(e,f),n(e,m),r(l,c,s),r(l,p,s),n(p,y),n(y,S),n(y,T),n(T,H),n(y,D),n(p,E),n(p,P),n(p,I),n(p,j),r(l,B,s),he($,l,s),r(l,N,s),r(l,q,s),r(l,g,s),r(l,b,s),n(b,O),n(b,z),n(b,F),n(F,h),n(h,C),n(h,ee),n(ee,ve),n(h,je),n(b,De),A&&A.m(b,null),r(l,ge,s),r(l,se,s),r(l,we,s),r(l,W,s),n(W,$e),n(W,Ne),n(W,U),n(U,Ce),n(U,Ve),L&&L.m(U,null),n(U,Se);for(let k=0;k<V.length;k+=1)V[k]&&V[k].m(U,null);r(l,Te,s),r(l,ie,s),r(l,qe,s),r(l,Y,s),n(Y,Oe),n(Y,Ee),n(Y,oe),n(oe,G),n(G,Me),n(G,Ie),n(G,He),n(G,Ue),n(G,M),n(M,Qe),he(te,M,null),n(M,ze),n(M,Ke),n(M,We),n(M,Le),n(M,Ye),n(M,Pe),n(M,Ge),n(M,Xe),n(M,Ze),n(M,Fe),n(M,xe),n(oe,et),he(le,oe,null),r(l,Re,s),r(l,re,s),r(l,Ae,s),r(l,X,s),n(X,de);for(let k=0;k<J.length;k+=1)J[k]&&J[k].m(de,null);n(X,lt),n(X,ce);for(let k=0;k<R.length;k+=1)R[k]&&R[k].m(ce,null);Z=!0},p(l,[s]){var ut,ft,mt,bt,_t;(!Z||s&1)&&a!==(a=l[0].name+"")&&x(f,a),(!Z||s&1)&&w!==(w=l[0].name+"")&&x(H,w);const k={};s&51&&(k.js=`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${l[5]}');

...

// example create data
const data = ${JSON.stringify(Object.assign({},l[4],Q.dummyCollectionSchemaData(l[0])),null,4)};

const record = await pb.collection('${(ut=l[0])==null?void 0:ut.name}').create(data);
`+(l[1]?`
// (optional) send an email verification request
await pb.collection('${(ft=l[0])==null?void 0:ft.name}').requestVerification('test@example.com');
`:"")),s&51&&(k.dart=`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${l[5]}');

...

// example create body
final body = <String, dynamic>${JSON.stringify(Object.assign({},l[4],Q.dummyCollectionSchemaData(l[0])),null,2)};

final record = await pb.collection('${(mt=l[0])==null?void 0:mt.name}').create(body: body);
`+(l[1]?`
// (optional) send an email verification request
await pb.collection('${(bt=l[0])==null?void 0:bt.name}').requestVerification('test@example.com');
`:"")),$.$set(k),(!Z||s&1)&&K!==(K=l[0].name+"")&&x(ve,K),l[6]?A||(A=gt(),A.c(),A.m(b,null)):A&&(A.d(1),A=null),l[1]?L?L.p(l,s):(L=wt(l),L.c(),L.m(U,Se)):L&&(L.d(1),L=null),s&1&&(me=ne((_t=l[0])==null?void 0:_t.schema),V=Be(V,s,at,1,l,me,Je,U,ht,$t,null,vt)),s&12&&(be=ne(l[3]),J=Be(J,s,st,1,l,be,tt,de,ht,Ct,null,yt)),s&12&&(pe=ne(l[3]),Ht(),R=Be(R,s,it,1,l,pe,nt,ce,Lt,St,null,kt),Pt())},i(l){if(!Z){ue($.$$.fragment,l),ue(te.$$.fragment,l),ue(le.$$.fragment,l);for(let s=0;s<pe.length;s+=1)ue(R[s]);Z=!0}},o(l){fe($.$$.fragment,l),fe(te.$$.fragment,l),fe(le.$$.fragment,l);for(let s=0;s<R.length;s+=1)fe(R[s]);Z=!1},d(l){l&&(d(e),d(c),d(p),d(B),d(N),d(q),d(g),d(b),d(ge),d(se),d(we),d(W),d(Te),d(ie),d(qe),d(Y),d(Re),d(re),d(Ae),d(X)),ke($,l),A&&A.d(),L&&L.d();for(let s=0;s<V.length;s+=1)V[s].d();ke(te),ke(le);for(let s=0;s<J.length;s+=1)J[s].d();for(let s=0;s<R.length;s+=1)R[s].d()}}}function Yt(o,e,t){let a,f,m,{collection:c}=e,p=200,y=[],S={};const T=w=>t(2,p=w.code);return o.$$set=w=>{"collection"in w&&t(0,c=w.collection)},o.$$.update=()=>{var w,H;o.$$.dirty&1&&t(1,a=c.type==="auth"),o.$$.dirty&1&&t(6,f=(c==null?void 0:c.createRule)===null),o.$$.dirty&1&&t(3,y=[{code:200,body:JSON.stringify(Q.dummyCollectionRecord(c),null,2)},{code:400,body:`
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
            `}]),o.$$.dirty&2&&(a?t(4,S={username:"test_username",email:"test@example.com",emailVisibility:!0,password:"12345678",passwordConfirm:"12345678"}):t(4,S={}))},t(5,m=Q.getApiExampleUrl(Ft.baseUrl)),[c,a,p,y,S,m,f,T]}class xt extends qt{constructor(e){super(),Ot(this,e,Yt,Wt,Mt,{collection:0})}}export{xt as default};
