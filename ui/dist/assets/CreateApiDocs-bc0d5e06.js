import{S as Pt,i as Ft,s as Rt,C as Q,N as ne,M as Lt,e as i,w as _,b as f,c as _e,f as v,g as r,h as n,m as he,x,P as Be,Q as $t,k as Bt,R as jt,n as Dt,t as ue,a as fe,o as d,d as ke,U as Nt,p as Vt,r as ye,u as Jt,y as ae}from"./index-a084d9d7.js";import{S as Et}from"./SdkTabs-ba0ec979.js";import{F as It}from"./FieldsQueryParam-71e01e64.js";function Ct(o,e,t){const a=o.slice();return a[7]=e[t],a}function St(o,e,t){const a=o.slice();return a[7]=e[t],a}function Tt(o,e,t){const a=o.slice();return a[12]=e[t],a}function qt(o){let e;return{c(){e=i("p"),e.innerHTML="Requires admin <code>Authorization:TOKEN</code> header",v(e,"class","txt-hint txt-sm txt-right")},m(t,a){r(t,e,a)},d(t){t&&d(e)}}}function Mt(o){let e,t,a,m,u,p,c,y,S,g,q,B,D,E,H,I,j,R,$,N,T,w,b;function M(h,C){var ee,K;return(K=(ee=h[0])==null?void 0:ee.options)!=null&&K.requireEmail?Qt:Ut}let z=M(o),L=z(o);return{c(){e=i("tr"),e.innerHTML='<td colspan="3" class="txt-hint">Auth fields</td>',t=f(),a=i("tr"),a.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>username</span></div></td> <td><span class="label">String</span></td> <td>The username of the auth record.
                    <br/>
                    If not set, it will be auto generated.</td>`,m=f(),u=i("tr"),p=i("td"),c=i("div"),L.c(),y=f(),S=i("span"),S.textContent="email",g=f(),q=i("td"),q.innerHTML='<span class="label">String</span>',B=f(),D=i("td"),D.textContent="Auth record email address.",E=f(),H=i("tr"),H.innerHTML='<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>emailVisibility</span></div></td> <td><span class="label">Boolean</span></td> <td>Whether to show/hide the auth record email when fetching the record data.</td>',I=f(),j=i("tr"),j.innerHTML='<td><div class="inline-flex"><span class="label label-success">Required</span> <span>password</span></div></td> <td><span class="label">String</span></td> <td>Auth record password.</td>',R=f(),$=i("tr"),$.innerHTML='<td><div class="inline-flex"><span class="label label-success">Required</span> <span>passwordConfirm</span></div></td> <td><span class="label">String</span></td> <td>Auth record password confirmation.</td>',N=f(),T=i("tr"),T.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>verified</span></div></td> <td><span class="label">Boolean</span></td> <td>Indicates whether the auth record is verified or not.
                    <br/>
                    This field can be set only by admins or auth records with &quot;Manage&quot; access.</td>`,w=f(),b=i("tr"),b.innerHTML='<td colspan="3" class="txt-hint">Schema fields</td>',v(c,"class","inline-flex")},m(h,C){r(h,e,C),r(h,t,C),r(h,a,C),r(h,m,C),r(h,u,C),n(u,p),n(p,c),L.m(c,null),n(c,y),n(c,S),n(u,g),n(u,q),n(u,B),n(u,D),r(h,E,C),r(h,H,C),r(h,I,C),r(h,j,C),r(h,R,C),r(h,$,C),r(h,N,C),r(h,T,C),r(h,w,C),r(h,b,C)},p(h,C){z!==(z=M(h))&&(L.d(1),L=z(h),L&&(L.c(),L.m(c,y)))},d(h){h&&(d(e),d(t),d(a),d(m),d(u),d(E),d(H),d(I),d(j),d(R),d($),d(N),d(T),d(w),d(b)),L.d()}}}function Ut(o){let e;return{c(){e=i("span"),e.textContent="Optional",v(e,"class","label label-warning")},m(t,a){r(t,e,a)},d(t){t&&d(e)}}}function Qt(o){let e;return{c(){e=i("span"),e.textContent="Required",v(e,"class","label label-success")},m(t,a){r(t,e,a)},d(t){t&&d(e)}}}function zt(o){let e;return{c(){e=i("span"),e.textContent="Optional",v(e,"class","label label-warning")},m(t,a){r(t,e,a)},d(t){t&&d(e)}}}function Kt(o){let e;return{c(){e=i("span"),e.textContent="Required",v(e,"class","label label-success")},m(t,a){r(t,e,a)},d(t){t&&d(e)}}}function Wt(o){var u;let e,t=((u=o[12].options)==null?void 0:u.maxSelect)===1?"id":"ids",a,m;return{c(){e=_("Relation record "),a=_(t),m=_(".")},m(p,c){r(p,e,c),r(p,a,c),r(p,m,c)},p(p,c){var y;c&1&&t!==(t=((y=p[12].options)==null?void 0:y.maxSelect)===1?"id":"ids")&&x(a,t)},d(p){p&&(d(e),d(a),d(m))}}}function Yt(o){let e,t,a,m,u;return{c(){e=_("File object."),t=i("br"),a=_(`
                        Set to `),m=i("code"),m.textContent="null",u=_(" to delete already uploaded file(s).")},m(p,c){r(p,e,c),r(p,t,c),r(p,a,c),r(p,m,c),r(p,u,c)},p:ae,d(p){p&&(d(e),d(t),d(a),d(m),d(u))}}}function Gt(o){let e;return{c(){e=_("URL address.")},m(t,a){r(t,e,a)},p:ae,d(t){t&&d(e)}}}function Xt(o){let e;return{c(){e=_("Email address.")},m(t,a){r(t,e,a)},p:ae,d(t){t&&d(e)}}}function Zt(o){let e;return{c(){e=_("JSON array or object.")},m(t,a){r(t,e,a)},p:ae,d(t){t&&d(e)}}}function xt(o){let e;return{c(){e=_("Number value.")},m(t,a){r(t,e,a)},p:ae,d(t){t&&d(e)}}}function el(o){let e;return{c(){e=_("Plain text value.")},m(t,a){r(t,e,a)},p:ae,d(t){t&&d(e)}}}function Ot(o,e){let t,a,m,u,p,c=e[12].name+"",y,S,g,q,B=Q.getFieldValueType(e[12])+"",D,E,H,I;function j(b,M){return b[12].required?Kt:zt}let R=j(e),$=R(e);function N(b,M){if(b[12].type==="text")return el;if(b[12].type==="number")return xt;if(b[12].type==="json")return Zt;if(b[12].type==="email")return Xt;if(b[12].type==="url")return Gt;if(b[12].type==="file")return Yt;if(b[12].type==="relation")return Wt}let T=N(e),w=T&&T(e);return{key:o,first:null,c(){t=i("tr"),a=i("td"),m=i("div"),$.c(),u=f(),p=i("span"),y=_(c),S=f(),g=i("td"),q=i("span"),D=_(B),E=f(),H=i("td"),w&&w.c(),I=f(),v(m,"class","inline-flex"),v(q,"class","label"),this.first=t},m(b,M){r(b,t,M),n(t,a),n(a,m),$.m(m,null),n(m,u),n(m,p),n(p,y),n(t,S),n(t,g),n(g,q),n(q,D),n(t,E),n(t,H),w&&w.m(H,null),n(t,I)},p(b,M){e=b,R!==(R=j(e))&&($.d(1),$=R(e),$&&($.c(),$.m(m,u))),M&1&&c!==(c=e[12].name+"")&&x(y,c),M&1&&B!==(B=Q.getFieldValueType(e[12])+"")&&x(D,B),T===(T=N(e))&&w?w.p(e,M):(w&&w.d(1),w=T&&T(e),w&&(w.c(),w.m(H,null)))},d(b){b&&d(t),$.d(),w&&w.d()}}}function At(o,e){let t,a=e[7].code+"",m,u,p,c;function y(){return e[6](e[7])}return{key:o,first:null,c(){t=i("button"),m=_(a),u=f(),v(t,"class","tab-item"),ye(t,"active",e[1]===e[7].code),this.first=t},m(S,g){r(S,t,g),n(t,m),n(t,u),p||(c=Jt(t,"click",y),p=!0)},p(S,g){e=S,g&4&&a!==(a=e[7].code+"")&&x(m,a),g&6&&ye(t,"active",e[1]===e[7].code)},d(S){S&&d(t),p=!1,c()}}}function Ht(o,e){let t,a,m,u;return a=new Lt({props:{content:e[7].body}}),{key:o,first:null,c(){t=i("div"),_e(a.$$.fragment),m=f(),v(t,"class","tab-item"),ye(t,"active",e[1]===e[7].code),this.first=t},m(p,c){r(p,t,c),he(a,t,null),n(t,m),u=!0},p(p,c){e=p;const y={};c&4&&(y.content=e[7].body),a.$set(y),(!u||c&6)&&ye(t,"active",e[1]===e[7].code)},i(p){u||(ue(a.$$.fragment,p),u=!0)},o(p){fe(a.$$.fragment,p),u=!1},d(p){p&&d(t),ke(a)}}}function tl(o){var ot,rt,dt,ct,pt,ut,ft,mt;let e,t,a=o[0].name+"",m,u,p,c,y,S,g,q=o[0].name+"",B,D,E,H,I,j,R,$,N,T,w,b,M,z,L,h,C,ee,K=o[0].name+"",ve,je,De,ge,se,we,W,$e,Ne,U,Ce,Ve,Se,V=[],Je=new Map,Te,ie,qe,Y,Me,Ee,oe,G,Oe,Ie,Ae,Ue,O,Qe,te,ze,Ke,We,He,Ye,Le,Ge,Xe,Ze,Pe,xe,et,le,Fe,re,Re,X,de,J=[],tt=new Map,lt,ce,P=[],nt=new Map,Z;$=new Et({props:{js:`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${o[4]}');

...

// example create data
const data = ${JSON.stringify(Object.assign({},o[3],Q.dummyCollectionSchemaData(o[0])),null,4)};

const record = await pb.collection('${(ot=o[0])==null?void 0:ot.name}').create(data);
`+((rt=o[0])!=null&&rt.isAuth?`
// (optional) send an email verification request
await pb.collection('${(dt=o[0])==null?void 0:dt.name}').requestVerification('test@example.com');
`:""),dart:`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${o[4]}');

...

// example create body
final body = <String, dynamic>${JSON.stringify(Object.assign({},o[3],Q.dummyCollectionSchemaData(o[0])),null,2)};

final record = await pb.collection('${(ct=o[0])==null?void 0:ct.name}').create(body: body);
`+((pt=o[0])!=null&&pt.isAuth?`
// (optional) send an email verification request
await pb.collection('${(ut=o[0])==null?void 0:ut.name}').requestVerification('test@example.com');
`:"")}});let F=o[5]&&qt(),A=((ft=o[0])==null?void 0:ft.isAuth)&&Mt(o),me=ne((mt=o[0])==null?void 0:mt.schema);const at=l=>l[12].name;for(let l=0;l<me.length;l+=1){let s=Tt(o,me,l),k=at(s);Je.set(k,V[l]=Ot(k,s))}te=new Lt({props:{content:"?expand=relField1,relField2.subRelField"}}),le=new It({});let be=ne(o[2]);const st=l=>l[7].code;for(let l=0;l<be.length;l+=1){let s=St(o,be,l),k=st(s);tt.set(k,J[l]=At(k,s))}let pe=ne(o[2]);const it=l=>l[7].code;for(let l=0;l<pe.length;l+=1){let s=Ct(o,pe,l),k=it(s);nt.set(k,P[l]=Ht(k,s))}return{c(){e=i("h3"),t=_("Create ("),m=_(a),u=_(")"),p=f(),c=i("div"),y=i("p"),S=_("Create a new "),g=i("strong"),B=_(q),D=_(" record."),E=f(),H=i("p"),H.innerHTML=`Body parameters could be sent as <code>application/json</code> or
        <code>multipart/form-data</code>.`,I=f(),j=i("p"),j.innerHTML=`File upload is supported only via <code>multipart/form-data</code>.
        <br/>
        For more info and examples you could check the detailed
        <a href="https://pocketbase.io/docs/files-handling/" target="_blank" rel="noopener noreferrer">Files upload and handling docs
        </a>.`,R=f(),_e($.$$.fragment),N=f(),T=i("h6"),T.textContent="API details",w=f(),b=i("div"),M=i("strong"),M.textContent="POST",z=f(),L=i("div"),h=i("p"),C=_("/api/collections/"),ee=i("strong"),ve=_(K),je=_("/records"),De=f(),F&&F.c(),ge=f(),se=i("div"),se.textContent="Body Parameters",we=f(),W=i("table"),$e=i("thead"),$e.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr>',Ne=f(),U=i("tbody"),Ce=i("tr"),Ce.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>id</span></div></td> <td><span class="label">String</span></td> <td><strong>15 characters string</strong> to store as record ID.
                <br/>
                If not set, it will be auto generated.</td>`,Ve=f(),A&&A.c(),Se=f();for(let l=0;l<V.length;l+=1)V[l].c();Te=f(),ie=i("div"),ie.textContent="Query parameters",qe=f(),Y=i("table"),Me=i("thead"),Me.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr>',Ee=f(),oe=i("tbody"),G=i("tr"),Oe=i("td"),Oe.textContent="expand",Ie=f(),Ae=i("td"),Ae.innerHTML='<span class="label">String</span>',Ue=f(),O=i("td"),Qe=_(`Auto expand relations when returning the created record. Ex.:
                `),_e(te.$$.fragment),ze=_(`
                Supports up to 6-levels depth nested relations expansion. `),Ke=i("br"),We=_(`
                The expanded relations will be appended to the record under the
                `),He=i("code"),He.textContent="expand",Ye=_(" property (eg. "),Le=i("code"),Le.textContent='"expand": {"relField1": {...}, ...}',Ge=_(`).
                `),Xe=i("br"),Ze=_(`
                Only the relations to which the request user has permissions to `),Pe=i("strong"),Pe.textContent="view",xe=_(" will be expanded."),et=f(),_e(le.$$.fragment),Fe=f(),re=i("div"),re.textContent="Responses",Re=f(),X=i("div"),de=i("div");for(let l=0;l<J.length;l+=1)J[l].c();lt=f(),ce=i("div");for(let l=0;l<P.length;l+=1)P[l].c();v(e,"class","m-b-sm"),v(c,"class","content txt-lg m-b-sm"),v(T,"class","m-b-xs"),v(M,"class","label label-primary"),v(L,"class","content"),v(b,"class","alert alert-success"),v(se,"class","section-title"),v(W,"class","table-compact table-border m-b-base"),v(ie,"class","section-title"),v(Y,"class","table-compact table-border m-b-base"),v(re,"class","section-title"),v(de,"class","tabs-header compact left"),v(ce,"class","tabs-content"),v(X,"class","tabs")},m(l,s){r(l,e,s),n(e,t),n(e,m),n(e,u),r(l,p,s),r(l,c,s),n(c,y),n(y,S),n(y,g),n(g,B),n(y,D),n(c,E),n(c,H),n(c,I),n(c,j),r(l,R,s),he($,l,s),r(l,N,s),r(l,T,s),r(l,w,s),r(l,b,s),n(b,M),n(b,z),n(b,L),n(L,h),n(h,C),n(h,ee),n(ee,ve),n(h,je),n(b,De),F&&F.m(b,null),r(l,ge,s),r(l,se,s),r(l,we,s),r(l,W,s),n(W,$e),n(W,Ne),n(W,U),n(U,Ce),n(U,Ve),A&&A.m(U,null),n(U,Se);for(let k=0;k<V.length;k+=1)V[k]&&V[k].m(U,null);r(l,Te,s),r(l,ie,s),r(l,qe,s),r(l,Y,s),n(Y,Me),n(Y,Ee),n(Y,oe),n(oe,G),n(G,Oe),n(G,Ie),n(G,Ae),n(G,Ue),n(G,O),n(O,Qe),he(te,O,null),n(O,ze),n(O,Ke),n(O,We),n(O,He),n(O,Ye),n(O,Le),n(O,Ge),n(O,Xe),n(O,Ze),n(O,Pe),n(O,xe),n(oe,et),he(le,oe,null),r(l,Fe,s),r(l,re,s),r(l,Re,s),r(l,X,s),n(X,de);for(let k=0;k<J.length;k+=1)J[k]&&J[k].m(de,null);n(X,lt),n(X,ce);for(let k=0;k<P.length;k+=1)P[k]&&P[k].m(ce,null);Z=!0},p(l,[s]){var bt,_t,ht,kt,yt,vt,gt,wt;(!Z||s&1)&&a!==(a=l[0].name+"")&&x(m,a),(!Z||s&1)&&q!==(q=l[0].name+"")&&x(B,q);const k={};s&25&&(k.js=`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${l[4]}');

...

// example create data
const data = ${JSON.stringify(Object.assign({},l[3],Q.dummyCollectionSchemaData(l[0])),null,4)};

const record = await pb.collection('${(bt=l[0])==null?void 0:bt.name}').create(data);
`+((_t=l[0])!=null&&_t.isAuth?`
// (optional) send an email verification request
await pb.collection('${(ht=l[0])==null?void 0:ht.name}').requestVerification('test@example.com');
`:"")),s&25&&(k.dart=`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${l[4]}');

...

// example create body
final body = <String, dynamic>${JSON.stringify(Object.assign({},l[3],Q.dummyCollectionSchemaData(l[0])),null,2)};

final record = await pb.collection('${(kt=l[0])==null?void 0:kt.name}').create(body: body);
`+((yt=l[0])!=null&&yt.isAuth?`
// (optional) send an email verification request
await pb.collection('${(vt=l[0])==null?void 0:vt.name}').requestVerification('test@example.com');
`:"")),$.$set(k),(!Z||s&1)&&K!==(K=l[0].name+"")&&x(ve,K),l[5]?F||(F=qt(),F.c(),F.m(b,null)):F&&(F.d(1),F=null),(gt=l[0])!=null&&gt.isAuth?A?A.p(l,s):(A=Mt(l),A.c(),A.m(U,Se)):A&&(A.d(1),A=null),s&1&&(me=ne((wt=l[0])==null?void 0:wt.schema),V=Be(V,s,at,1,l,me,Je,U,$t,Ot,null,Tt)),s&6&&(be=ne(l[2]),J=Be(J,s,st,1,l,be,tt,de,$t,At,null,St)),s&6&&(pe=ne(l[2]),Bt(),P=Be(P,s,it,1,l,pe,nt,ce,jt,Ht,null,Ct),Dt())},i(l){if(!Z){ue($.$$.fragment,l),ue(te.$$.fragment,l),ue(le.$$.fragment,l);for(let s=0;s<pe.length;s+=1)ue(P[s]);Z=!0}},o(l){fe($.$$.fragment,l),fe(te.$$.fragment,l),fe(le.$$.fragment,l);for(let s=0;s<P.length;s+=1)fe(P[s]);Z=!1},d(l){l&&(d(e),d(p),d(c),d(R),d(N),d(T),d(w),d(b),d(ge),d(se),d(we),d(W),d(Te),d(ie),d(qe),d(Y),d(Fe),d(re),d(Re),d(X)),ke($,l),F&&F.d(),A&&A.d();for(let s=0;s<V.length;s+=1)V[s].d();ke(te),ke(le);for(let s=0;s<J.length;s+=1)J[s].d();for(let s=0;s<P.length;s+=1)P[s].d()}}}function ll(o,e,t){let a,m,{collection:u=new Nt}=e,p=200,c=[],y={};const S=g=>t(1,p=g.code);return o.$$set=g=>{"collection"in g&&t(0,u=g.collection)},o.$$.update=()=>{var g,q;o.$$.dirty&1&&t(5,a=(u==null?void 0:u.createRule)===null),o.$$.dirty&1&&t(2,c=[{code:200,body:JSON.stringify(Q.dummyCollectionRecord(u),null,2)},{code:400,body:`
                {
                  "code": 400,
                  "message": "Failed to create record.",
                  "data": {
                    "${(q=(g=u==null?void 0:u.schema)==null?void 0:g[0])==null?void 0:q.name}": {
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
            `}]),o.$$.dirty&1&&(u.$isAuth?t(3,y={username:"test_username",email:"test@example.com",emailVisibility:!0,password:"12345678",passwordConfirm:"12345678"}):t(3,y={}))},t(4,m=Q.getApiExampleUrl(Vt.baseUrl)),[u,p,c,y,m,a,S]}class il extends Pt{constructor(e){super(),Ft(this,e,ll,tl,Rt,{collection:0})}}export{il as default};
