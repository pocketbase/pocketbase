import{S as $t,i as Mt,s as St,V as Ot,X as se,W as Ct,h as o,d as ge,t as _e,a as he,I as ee,Z as Ue,_ as bt,C as qt,$ as Rt,D as Ht,l as r,n,m as we,u as i,A as h,v as f,c as Te,w as k,J as ye,p as Lt,k as Ce,o as Pt,H as ie}from"./index-0unWA3Bg.js";import{F as Dt}from"./FieldsQueryParam-ZLlGrCBp.js";function mt(d,e,t){const a=d.slice();return a[10]=e[t],a}function _t(d,e,t){const a=d.slice();return a[10]=e[t],a}function ht(d,e,t){const a=d.slice();return a[15]=e[t],a}function yt(d){let e;return{c(){e=i("p"),e.innerHTML=`<em>Note that in case of a password change all previously issued tokens for the current record
                will be automatically invalidated and if you want your user to remain signed in you need to
                reauthenticate manually after the update call.</em>`},m(t,a){r(t,e,a)},d(t){t&&o(e)}}}function kt(d){let e;return{c(){e=i("p"),e.innerHTML="Requires superuser <code>Authorization:TOKEN</code> header",k(e,"class","txt-hint txt-sm txt-right")},m(t,a){r(t,e,a)},d(t){t&&o(e)}}}function vt(d){let e,t,a,m,p,c,u,b,S,C,H,L,$,M,q,P,U,J,O,R,D,v,g,w;function x(_,T){var te,W,le;return T&1&&(b=null),b==null&&(b=!!((le=(W=(te=_[0])==null?void 0:te.fields)==null?void 0:W.find(Wt))!=null&&le.required)),b?Bt:Ft}let Q=x(d,-1),B=Q(d);return{c(){e=i("tr"),e.innerHTML='<td colspan="3" class="txt-hint txt-bold">Auth specific fields</td>',t=f(),a=i("tr"),a.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>email</span></div></td> <td><span class="label">String</span></td> <td>The auth record email address.
                    <br/>
                    This field can be updated only by superusers or auth records with &quot;Manage&quot; access.
                    <br/>
                    Regular accounts can update their email by calling &quot;Request email change&quot;.</td>`,m=f(),p=i("tr"),c=i("td"),u=i("div"),B.c(),S=f(),C=i("span"),C.textContent="emailVisibility",H=f(),L=i("td"),L.innerHTML='<span class="label">Boolean</span>',$=f(),M=i("td"),M.textContent="Whether to show/hide the auth record email when fetching the record data.",q=f(),P=i("tr"),P.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>oldPassword</span></div></td> <td><span class="label">String</span></td> <td>Old auth record password.
                    <br/>
                    This field is required only when changing the record password. Superusers and auth records
                    with &quot;Manage&quot; access can skip this field.</td>`,U=f(),J=i("tr"),J.innerHTML='<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>password</span></div></td> <td><span class="label">String</span></td> <td>New auth record password.</td>',O=f(),R=i("tr"),R.innerHTML='<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>passwordConfirm</span></div></td> <td><span class="label">String</span></td> <td>New auth record password confirmation.</td>',D=f(),v=i("tr"),v.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>verified</span></div></td> <td><span class="label">Boolean</span></td> <td>Indicates whether the auth record is verified or not.
                    <br/>
                    This field can be set only by superusers or auth records with &quot;Manage&quot; access.</td>`,g=f(),w=i("tr"),w.innerHTML='<td colspan="3" class="txt-hint txt-bold">Other fields</td>',k(u,"class","inline-flex")},m(_,T){r(_,e,T),r(_,t,T),r(_,a,T),r(_,m,T),r(_,p,T),n(p,c),n(c,u),B.m(u,null),n(u,S),n(u,C),n(p,H),n(p,L),n(p,$),n(p,M),r(_,q,T),r(_,P,T),r(_,U,T),r(_,J,T),r(_,O,T),r(_,R,T),r(_,D,T),r(_,v,T),r(_,g,T),r(_,w,T)},p(_,T){Q!==(Q=x(_,T))&&(B.d(1),B=Q(_),B&&(B.c(),B.m(u,S)))},d(_){_&&(o(e),o(t),o(a),o(m),o(p),o(q),o(P),o(U),o(J),o(O),o(R),o(D),o(v),o(g),o(w)),B.d()}}}function Ft(d){let e;return{c(){e=i("span"),e.textContent="Optional",k(e,"class","label label-warning")},m(t,a){r(t,e,a)},d(t){t&&o(e)}}}function Bt(d){let e;return{c(){e=i("span"),e.textContent="Required",k(e,"class","label label-success")},m(t,a){r(t,e,a)},d(t){t&&o(e)}}}function Nt(d){let e;return{c(){e=i("span"),e.textContent="Optional",k(e,"class","label label-warning")},m(t,a){r(t,e,a)},d(t){t&&o(e)}}}function At(d){let e;return{c(){e=i("span"),e.textContent="Required",k(e,"class","label label-success")},m(t,a){r(t,e,a)},d(t){t&&o(e)}}}function Et(d){let e,t=d[15].maxSelect==1?"id":"ids",a,m;return{c(){e=h("Relation record "),a=h(t),m=h(".")},m(p,c){r(p,e,c),r(p,a,c),r(p,m,c)},p(p,c){c&32&&t!==(t=p[15].maxSelect==1?"id":"ids")&&ee(a,t)},d(p){p&&(o(e),o(a),o(m))}}}function It(d){let e,t,a,m,p;return{c(){e=h("File object."),t=i("br"),a=h(`
                        Set to `),m=i("code"),m.textContent="null",p=h(" to delete already uploaded file(s).")},m(c,u){r(c,e,u),r(c,t,u),r(c,a,u),r(c,m,u),r(c,p,u)},p:ie,d(c){c&&(o(e),o(t),o(a),o(m),o(p))}}}function Jt(d){let e;return{c(){e=h("URL address.")},m(t,a){r(t,e,a)},p:ie,d(t){t&&o(e)}}}function Ut(d){let e;return{c(){e=h("Email address.")},m(t,a){r(t,e,a)},p:ie,d(t){t&&o(e)}}}function jt(d){let e;return{c(){e=h("JSON array or object.")},m(t,a){r(t,e,a)},p:ie,d(t){t&&o(e)}}}function Vt(d){let e;return{c(){e=h("Number value.")},m(t,a){r(t,e,a)},p:ie,d(t){t&&o(e)}}}function xt(d){let e;return{c(){e=h("Plain text value.")},m(t,a){r(t,e,a)},p:ie,d(t){t&&o(e)}}}function gt(d,e){let t,a,m,p,c,u=e[15].name+"",b,S,C,H,L=ye.getFieldValueType(e[15])+"",$,M,q,P;function U(g,w){return g[15].required?At:Nt}let J=U(e),O=J(e);function R(g,w){if(g[15].type==="text")return xt;if(g[15].type==="number")return Vt;if(g[15].type==="json")return jt;if(g[15].type==="email")return Ut;if(g[15].type==="url")return Jt;if(g[15].type==="file")return It;if(g[15].type==="relation")return Et}let D=R(e),v=D&&D(e);return{key:d,first:null,c(){t=i("tr"),a=i("td"),m=i("div"),O.c(),p=f(),c=i("span"),b=h(u),S=f(),C=i("td"),H=i("span"),$=h(L),M=f(),q=i("td"),v&&v.c(),P=f(),k(m,"class","inline-flex"),k(H,"class","label"),this.first=t},m(g,w){r(g,t,w),n(t,a),n(a,m),O.m(m,null),n(m,p),n(m,c),n(c,b),n(t,S),n(t,C),n(C,H),n(H,$),n(t,M),n(t,q),v&&v.m(q,null),n(t,P)},p(g,w){e=g,J!==(J=U(e))&&(O.d(1),O=J(e),O&&(O.c(),O.m(m,p))),w&32&&u!==(u=e[15].name+"")&&ee(b,u),w&32&&L!==(L=ye.getFieldValueType(e[15])+"")&&ee($,L),D===(D=R(e))&&v?v.p(e,w):(v&&v.d(1),v=D&&D(e),v&&(v.c(),v.m(q,null)))},d(g){g&&o(t),O.d(),v&&v.d()}}}function wt(d,e){let t,a=e[10].code+"",m,p,c,u;function b(){return e[9](e[10])}return{key:d,first:null,c(){t=i("button"),m=h(a),p=f(),k(t,"class","tab-item"),Ce(t,"active",e[2]===e[10].code),this.first=t},m(S,C){r(S,t,C),n(t,m),n(t,p),c||(u=Pt(t,"click",b),c=!0)},p(S,C){e=S,C&8&&a!==(a=e[10].code+"")&&ee(m,a),C&12&&Ce(t,"active",e[2]===e[10].code)},d(S){S&&o(t),c=!1,u()}}}function Tt(d,e){let t,a,m,p;return a=new Ct({props:{content:e[10].body}}),{key:d,first:null,c(){t=i("div"),Te(a.$$.fragment),m=f(),k(t,"class","tab-item"),Ce(t,"active",e[2]===e[10].code),this.first=t},m(c,u){r(c,t,u),we(a,t,null),n(t,m),p=!0},p(c,u){e=c;const b={};u&8&&(b.content=e[10].body),a.$set(b),(!p||u&12)&&Ce(t,"active",e[2]===e[10].code)},i(c){p||(he(a.$$.fragment,c),p=!0)},o(c){_e(a.$$.fragment,c),p=!1},d(c){c&&o(t),ge(a)}}}function Qt(d){var ct,ut;let e,t,a=d[0].name+"",m,p,c,u,b,S,C,H=d[0].name+"",L,$,M,q,P,U,J,O,R,D,v,g,w,x,Q,B,_,T,te,W=d[0].name+"",le,je,$e,Ve,Me,de,Se,oe,Oe,re,qe,z,Re,xe,K,He,j=[],Qe=new Map,Le,ce,Pe,X,De,We,ue,Y,Fe,ze,Be,Ke,N,Xe,ne,Ye,Ze,Ge,Ne,et,Ae,tt,Ee,lt,nt,ae,Ie,pe,Je,Z,fe,V=[],at=new Map,st,be,A=[],it=new Map,G,E=d[1]&&yt();R=new Ot({props:{js:`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${d[4]}');

...

// example update data
const data = ${JSON.stringify(d[7](d[0]),null,4)};

const record = await pb.collection('${(ct=d[0])==null?void 0:ct.name}').update('RECORD_ID', data);
    `,dart:`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${d[4]}');

...

// example update body
final body = <String, dynamic>${JSON.stringify(d[7](d[0]),null,2)};

final record = await pb.collection('${(ut=d[0])==null?void 0:ut.name}').update('RECORD_ID', body: body);
    `}});let I=d[6]&&kt(),F=d[1]&&vt(d),ke=se(d[5]);const dt=l=>l[15].name;for(let l=0;l<ke.length;l+=1){let s=ht(d,ke,l),y=dt(s);Qe.set(y,j[l]=gt(y,s))}ne=new Ct({props:{content:"?expand=relField1,relField2.subRelField21"}}),ae=new Dt({});let ve=se(d[3]);const ot=l=>l[10].code;for(let l=0;l<ve.length;l+=1){let s=_t(d,ve,l),y=ot(s);at.set(y,V[l]=wt(y,s))}let me=se(d[3]);const rt=l=>l[10].code;for(let l=0;l<me.length;l+=1){let s=mt(d,me,l),y=rt(s);it.set(y,A[l]=Tt(y,s))}return{c(){e=i("h3"),t=h("Update ("),m=h(a),p=h(")"),c=f(),u=i("div"),b=i("p"),S=h("Update a single "),C=i("strong"),L=h(H),$=h(" record."),M=f(),q=i("p"),q.innerHTML=`Body parameters could be sent as <code>application/json</code> or
        <code>multipart/form-data</code>.`,P=f(),U=i("p"),U.innerHTML=`File upload is supported only via <code>multipart/form-data</code>.
        <br/>
        For more info and examples you could check the detailed
        <a href="https://pocketbase.io/docs/files-handling" target="_blank" rel="noopener noreferrer">Files upload and handling docs
        </a>.`,J=f(),E&&E.c(),O=f(),Te(R.$$.fragment),D=f(),v=i("h6"),v.textContent="API details",g=f(),w=i("div"),x=i("strong"),x.textContent="PATCH",Q=f(),B=i("div"),_=i("p"),T=h("/api/collections/"),te=i("strong"),le=h(W),je=h("/records/"),$e=i("strong"),$e.textContent=":id",Ve=f(),I&&I.c(),Me=f(),de=i("div"),de.textContent="Path parameters",Se=f(),oe=i("table"),oe.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr></thead> <tbody><tr><td>id</td> <td><span class="label">String</span></td> <td>ID of the record to update.</td></tr></tbody>',Oe=f(),re=i("div"),re.textContent="Body Parameters",qe=f(),z=i("table"),Re=i("thead"),Re.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr>',xe=f(),K=i("tbody"),F&&F.c(),He=f();for(let l=0;l<j.length;l+=1)j[l].c();Le=f(),ce=i("div"),ce.textContent="Query parameters",Pe=f(),X=i("table"),De=i("thead"),De.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr>',We=f(),ue=i("tbody"),Y=i("tr"),Fe=i("td"),Fe.textContent="expand",ze=f(),Be=i("td"),Be.innerHTML='<span class="label">String</span>',Ke=f(),N=i("td"),Xe=h(`Auto expand relations when returning the updated record. Ex.:
                `),Te(ne.$$.fragment),Ye=h(`
                Supports up to 6-levels depth nested relations expansion. `),Ze=i("br"),Ge=h(`
                The expanded relations will be appended to the record under the
                `),Ne=i("code"),Ne.textContent="expand",et=h(" property (eg. "),Ae=i("code"),Ae.textContent='"expand": {"relField1": {...}, ...}',tt=h(`). Only
                the relations that the user has permissions to `),Ee=i("strong"),Ee.textContent="view",lt=h(" will be expanded."),nt=f(),Te(ae.$$.fragment),Ie=f(),pe=i("div"),pe.textContent="Responses",Je=f(),Z=i("div"),fe=i("div");for(let l=0;l<V.length;l+=1)V[l].c();st=f(),be=i("div");for(let l=0;l<A.length;l+=1)A[l].c();k(e,"class","m-b-sm"),k(u,"class","content txt-lg m-b-sm"),k(v,"class","m-b-xs"),k(x,"class","label label-primary"),k(B,"class","content"),k(w,"class","alert alert-warning"),k(de,"class","section-title"),k(oe,"class","table-compact table-border m-b-base"),k(re,"class","section-title"),k(z,"class","table-compact table-border m-b-base"),k(ce,"class","section-title"),k(X,"class","table-compact table-border m-b-lg"),k(pe,"class","section-title"),k(fe,"class","tabs-header compact combined left"),k(be,"class","tabs-content"),k(Z,"class","tabs")},m(l,s){r(l,e,s),n(e,t),n(e,m),n(e,p),r(l,c,s),r(l,u,s),n(u,b),n(b,S),n(b,C),n(C,L),n(b,$),n(u,M),n(u,q),n(u,P),n(u,U),n(u,J),E&&E.m(u,null),r(l,O,s),we(R,l,s),r(l,D,s),r(l,v,s),r(l,g,s),r(l,w,s),n(w,x),n(w,Q),n(w,B),n(B,_),n(_,T),n(_,te),n(te,le),n(_,je),n(_,$e),n(w,Ve),I&&I.m(w,null),r(l,Me,s),r(l,de,s),r(l,Se,s),r(l,oe,s),r(l,Oe,s),r(l,re,s),r(l,qe,s),r(l,z,s),n(z,Re),n(z,xe),n(z,K),F&&F.m(K,null),n(K,He);for(let y=0;y<j.length;y+=1)j[y]&&j[y].m(K,null);r(l,Le,s),r(l,ce,s),r(l,Pe,s),r(l,X,s),n(X,De),n(X,We),n(X,ue),n(ue,Y),n(Y,Fe),n(Y,ze),n(Y,Be),n(Y,Ke),n(Y,N),n(N,Xe),we(ne,N,null),n(N,Ye),n(N,Ze),n(N,Ge),n(N,Ne),n(N,et),n(N,Ae),n(N,tt),n(N,Ee),n(N,lt),n(ue,nt),we(ae,ue,null),r(l,Ie,s),r(l,pe,s),r(l,Je,s),r(l,Z,s),n(Z,fe);for(let y=0;y<V.length;y+=1)V[y]&&V[y].m(fe,null);n(Z,st),n(Z,be);for(let y=0;y<A.length;y+=1)A[y]&&A[y].m(be,null);G=!0},p(l,[s]){var pt,ft;(!G||s&1)&&a!==(a=l[0].name+"")&&ee(m,a),(!G||s&1)&&H!==(H=l[0].name+"")&&ee(L,H),l[1]?E||(E=yt(),E.c(),E.m(u,null)):E&&(E.d(1),E=null);const y={};s&17&&(y.js=`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${l[4]}');

...

// example update data
const data = ${JSON.stringify(l[7](l[0]),null,4)};

const record = await pb.collection('${(pt=l[0])==null?void 0:pt.name}').update('RECORD_ID', data);
    `),s&17&&(y.dart=`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${l[4]}');

...

// example update body
final body = <String, dynamic>${JSON.stringify(l[7](l[0]),null,2)};

final record = await pb.collection('${(ft=l[0])==null?void 0:ft.name}').update('RECORD_ID', body: body);
    `),R.$set(y),(!G||s&1)&&W!==(W=l[0].name+"")&&ee(le,W),l[6]?I||(I=kt(),I.c(),I.m(w,null)):I&&(I.d(1),I=null),l[1]?F?F.p(l,s):(F=vt(l),F.c(),F.m(K,He)):F&&(F.d(1),F=null),s&32&&(ke=se(l[5]),j=Ue(j,s,dt,1,l,ke,Qe,K,bt,gt,null,ht)),s&12&&(ve=se(l[3]),V=Ue(V,s,ot,1,l,ve,at,fe,bt,wt,null,_t)),s&12&&(me=se(l[3]),qt(),A=Ue(A,s,rt,1,l,me,it,be,Rt,Tt,null,mt),Ht())},i(l){if(!G){he(R.$$.fragment,l),he(ne.$$.fragment,l),he(ae.$$.fragment,l);for(let s=0;s<me.length;s+=1)he(A[s]);G=!0}},o(l){_e(R.$$.fragment,l),_e(ne.$$.fragment,l),_e(ae.$$.fragment,l);for(let s=0;s<A.length;s+=1)_e(A[s]);G=!1},d(l){l&&(o(e),o(c),o(u),o(O),o(D),o(v),o(g),o(w),o(Me),o(de),o(Se),o(oe),o(Oe),o(re),o(qe),o(z),o(Le),o(ce),o(Pe),o(X),o(Ie),o(pe),o(Je),o(Z)),E&&E.d(),ge(R,l),I&&I.d(),F&&F.d();for(let s=0;s<j.length;s+=1)j[s].d();ge(ne),ge(ae);for(let s=0;s<V.length;s+=1)V[s].d();for(let s=0;s<A.length;s+=1)A[s].d()}}}const Wt=d=>d.name=="emailVisibility";function zt(d,e,t){let a,m,p,c,u,{collection:b}=e,S=200,C=[];function H($){let M=ye.dummyCollectionSchemaData($,!0);return a&&(M.oldPassword="12345678",M.password="87654321",M.passwordConfirm="87654321",delete M.verified,delete M.email),M}const L=$=>t(2,S=$.code);return d.$$set=$=>{"collection"in $&&t(0,b=$.collection)},d.$$.update=()=>{var $,M,q;d.$$.dirty&1&&t(1,a=(b==null?void 0:b.type)==="auth"),d.$$.dirty&1&&t(6,m=(b==null?void 0:b.updateRule)===null),d.$$.dirty&2&&t(8,p=a?["id","password","verified","email","emailVisibility"]:["id"]),d.$$.dirty&257&&t(5,c=(($=b==null?void 0:b.fields)==null?void 0:$.filter(P=>!P.hidden&&P.type!="autodate"&&!p.includes(P.name)))||[]),d.$$.dirty&1&&t(3,C=[{code:200,body:JSON.stringify(ye.dummyCollectionRecord(b),null,2)},{code:400,body:`
                {
                  "status": 400,
                  "message": "Failed to update record.",
                  "data": {
                    "${(q=(M=b==null?void 0:b.fields)==null?void 0:M[0])==null?void 0:q.name}": {
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
            `}])},t(4,u=ye.getApiExampleUrl(Lt.baseURL)),[b,a,S,C,u,c,m,H,p,L]}class Yt extends $t{constructor(e){super(),Mt(this,e,zt,Qt,St,{collection:0})}}export{Yt as default};
