import{S as Ot,i as St,s as Mt,V as $t,J as x,X as ie,W as Tt,h as i,z as h,j as f,c as ve,k,n as o,o as n,m as we,H as te,Y as Je,Z as bt,E as qt,_ as Rt,G as Ht,t as he,a as ye,v as r,d as Ce,p as Dt,l as Te,q as Lt,I as de}from"./index-SKn09NMF.js";import{F as Pt}from"./FieldsQueryParam-NXRpAlxi.js";function mt(d,e,t){const a=d.slice();return a[10]=e[t],a}function _t(d,e,t){const a=d.slice();return a[10]=e[t],a}function ht(d,e,t){const a=d.slice();return a[15]=e[t],a}function yt(d){let e;return{c(){e=i("p"),e.innerHTML=`<em>Note that in case of a password change all previously issued tokens for the current record
                will be automatically invalidated and if you want your user to remain signed in you need to
                reauthenticate manually after the update call.</em>`},m(t,a){o(t,e,a)},d(t){t&&r(e)}}}function kt(d){let e;return{c(){e=i("p"),e.innerHTML="Requires superuser <code>Authorization:TOKEN</code> header",k(e,"class","txt-hint txt-sm txt-right")},m(t,a){o(t,e,a)},d(t){t&&r(e)}}}function gt(d){let e,t,a,m,p,c,u,b,O,T,M,H,S,E,q,D,J,I,$,R,L,g,v,w;function z(_,C){var le,W,ne;return C&1&&(b=null),b==null&&(b=!!((ne=(W=(le=_[0])==null?void 0:le.fields)==null?void 0:W.find(Qt))!=null&&ne.required)),b?jt:Ft}let Q=z(d,-1),F=Q(d);return{c(){e=i("tr"),e.innerHTML='<td colspan="3" class="txt-hint txt-bold">Auth specific fields</td>',t=f(),a=i("tr"),a.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>email</span></div></td> <td><span class="label">String</span></td> <td>The auth record email address.
                    <br/>
                    This field can be updated only by superusers or auth records with &quot;Manage&quot; access.
                    <br/>
                    Regular accounts can update their email by calling &quot;Request email change&quot;.</td>`,m=f(),p=i("tr"),c=i("td"),u=i("div"),F.c(),O=f(),T=i("span"),T.textContent="emailVisibility",M=f(),H=i("td"),H.innerHTML='<span class="label">Boolean</span>',S=f(),E=i("td"),E.textContent="Whether to show/hide the auth record email when fetching the record data.",q=f(),D=i("tr"),D.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>oldPassword</span></div></td> <td><span class="label">String</span></td> <td>Old auth record password.
                    <br/>
                    This field is required only when changing the record password. Superusers and auth records
                    with &quot;Manage&quot; access can skip this field.</td>`,J=f(),I=i("tr"),I.innerHTML='<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>password</span></div></td> <td><span class="label">String</span></td> <td>New auth record password.</td>',$=f(),R=i("tr"),R.innerHTML='<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>passwordConfirm</span></div></td> <td><span class="label">String</span></td> <td>New auth record password confirmation.</td>',L=f(),g=i("tr"),g.innerHTML=`<td><div class="inline-flex"><span class="label label-warning">Optional</span> <span>verified</span></div></td> <td><span class="label">Boolean</span></td> <td>Indicates whether the auth record is verified or not.
                    <br/>
                    This field can be set only by superusers or auth records with &quot;Manage&quot; access.</td>`,v=f(),w=i("tr"),w.innerHTML='<td colspan="3" class="txt-hint txt-bold">Other fields</td>',k(u,"class","inline-flex")},m(_,C){o(_,e,C),o(_,t,C),o(_,a,C),o(_,m,C),o(_,p,C),n(p,c),n(c,u),F.m(u,null),n(u,O),n(u,T),n(p,M),n(p,H),n(p,S),n(p,E),o(_,q,C),o(_,D,C),o(_,J,C),o(_,I,C),o(_,$,C),o(_,R,C),o(_,L,C),o(_,g,C),o(_,v,C),o(_,w,C)},p(_,C){Q!==(Q=z(_,C))&&(F.d(1),F=Q(_),F&&(F.c(),F.m(u,O)))},d(_){_&&(r(e),r(t),r(a),r(m),r(p),r(q),r(D),r(J),r(I),r($),r(R),r(L),r(g),r(v),r(w)),F.d()}}}function Ft(d){let e;return{c(){e=i("span"),e.textContent="Optional",k(e,"class","label label-warning")},m(t,a){o(t,e,a)},d(t){t&&r(e)}}}function jt(d){let e;return{c(){e=i("span"),e.textContent="Required",k(e,"class","label label-success")},m(t,a){o(t,e,a)},d(t){t&&r(e)}}}function Bt(d){let e;return{c(){e=i("span"),e.textContent="Optional",k(e,"class","label label-warning")},m(t,a){o(t,e,a)},d(t){t&&r(e)}}}function Nt(d){let e;return{c(){e=i("span"),e.textContent="Required",k(e,"class","label label-success")},m(t,a){o(t,e,a)},d(t){t&&r(e)}}}function At(d){let e,t=d[15].maxSelect==1?"id":"ids",a,m;return{c(){e=h("Relation record "),a=h(t),m=h(".")},m(p,c){o(p,e,c),o(p,a,c),o(p,m,c)},p(p,c){c&64&&t!==(t=p[15].maxSelect==1?"id":"ids")&&te(a,t)},d(p){p&&(r(e),r(a),r(m))}}}function Et(d){let e,t,a,m,p;return{c(){e=h("File object."),t=i("br"),a=h(`
                        Set to `),m=i("code"),m.textContent="null",p=h(" to delete already uploaded file(s).")},m(c,u){o(c,e,u),o(c,t,u),o(c,a,u),o(c,m,u),o(c,p,u)},p:de,d(c){c&&(r(e),r(t),r(a),r(m),r(p))}}}function It(d){let e;return{c(){e=h("URL address.")},m(t,a){o(t,e,a)},p:de,d(t){t&&r(e)}}}function Jt(d){let e;return{c(){e=h("Email address.")},m(t,a){o(t,e,a)},p:de,d(t){t&&r(e)}}}function Ut(d){let e;return{c(){e=h("JSON array or object.")},m(t,a){o(t,e,a)},p:de,d(t){t&&r(e)}}}function Vt(d){let e;return{c(){e=h("Number value.")},m(t,a){o(t,e,a)},p:de,d(t){t&&r(e)}}}function xt(d){let e;return{c(){e=h("Plain text value.")},m(t,a){o(t,e,a)},p:de,d(t){t&&r(e)}}}function vt(d,e){let t,a,m,p,c,u=e[15].name+"",b,O,T,M,H=x.getFieldValueType(e[15])+"",S,E,q,D;function J(v,w){return v[15].required?Nt:Bt}let I=J(e),$=I(e);function R(v,w){if(v[15].type==="text")return xt;if(v[15].type==="number")return Vt;if(v[15].type==="json")return Ut;if(v[15].type==="email")return Jt;if(v[15].type==="url")return It;if(v[15].type==="file")return Et;if(v[15].type==="relation")return At}let L=R(e),g=L&&L(e);return{key:d,first:null,c(){t=i("tr"),a=i("td"),m=i("div"),$.c(),p=f(),c=i("span"),b=h(u),O=f(),T=i("td"),M=i("span"),S=h(H),E=f(),q=i("td"),g&&g.c(),D=f(),k(m,"class","inline-flex"),k(M,"class","label"),this.first=t},m(v,w){o(v,t,w),n(t,a),n(a,m),$.m(m,null),n(m,p),n(m,c),n(c,b),n(t,O),n(t,T),n(T,M),n(M,S),n(t,E),n(t,q),g&&g.m(q,null),n(t,D)},p(v,w){e=v,I!==(I=J(e))&&($.d(1),$=I(e),$&&($.c(),$.m(m,p))),w&64&&u!==(u=e[15].name+"")&&te(b,u),w&64&&H!==(H=x.getFieldValueType(e[15])+"")&&te(S,H),L===(L=R(e))&&g?g.p(e,w):(g&&g.d(1),g=L&&L(e),g&&(g.c(),g.m(q,null)))},d(v){v&&r(t),$.d(),g&&g.d()}}}function wt(d,e){let t,a=e[10].code+"",m,p,c,u;function b(){return e[9](e[10])}return{key:d,first:null,c(){t=i("button"),m=h(a),p=f(),k(t,"class","tab-item"),Te(t,"active",e[2]===e[10].code),this.first=t},m(O,T){o(O,t,T),n(t,m),n(t,p),c||(u=Lt(t,"click",b),c=!0)},p(O,T){e=O,T&8&&a!==(a=e[10].code+"")&&te(m,a),T&12&&Te(t,"active",e[2]===e[10].code)},d(O){O&&r(t),c=!1,u()}}}function Ct(d,e){let t,a,m,p;return a=new Tt({props:{content:e[10].body}}),{key:d,first:null,c(){t=i("div"),ve(a.$$.fragment),m=f(),k(t,"class","tab-item"),Te(t,"active",e[2]===e[10].code),this.first=t},m(c,u){o(c,t,u),we(a,t,null),n(t,m),p=!0},p(c,u){e=c;const b={};u&8&&(b.content=e[10].body),a.$set(b),(!p||u&12)&&Te(t,"active",e[2]===e[10].code)},i(c){p||(he(a.$$.fragment,c),p=!0)},o(c){ye(a.$$.fragment,c),p=!1},d(c){c&&r(t),Ce(a)}}}function zt(d){var ct,ut;let e,t,a=d[0].name+"",m,p,c,u,b,O,T,M=d[0].name+"",H,S,E,q,D,J,I,$,R,L,g,v,w,z,Q,F,_,C,le,W=d[0].name+"",ne,Ue,Oe,Ve,Se,oe,Me,re,$e,ce,qe,Y,Re,xe,G,He,U=[],ze=new Map,De,ue,Le,K,Pe,Qe,pe,X,Fe,We,je,Ye,j,Ge,ae,Ke,Xe,Ze,Be,et,Ne,tt,Ae,lt,nt,se,Ee,fe,Ie,Z,be,V=[],at=new Map,st,me,B=[],it=new Map,ee,N=d[1]&&yt();R=new $t({props:{js:`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${d[5]}');

...

// example update data
const data = ${JSON.stringify(Object.assign({},d[4],x.dummyCollectionSchemaData(d[0],!0)),null,4)};

const record = await pb.collection('${(ct=d[0])==null?void 0:ct.name}').update('RECORD_ID', data);
    `,dart:`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${d[5]}');

...

// example update body
final body = <String, dynamic>${JSON.stringify(Object.assign({},d[4],x.dummyCollectionSchemaData(d[0],!0)),null,2)};

final record = await pb.collection('${(ut=d[0])==null?void 0:ut.name}').update('RECORD_ID', body: body);
    `}});let A=d[7]&&kt(),P=d[1]&&gt(d),ke=ie(d[6]);const dt=l=>l[15].name;for(let l=0;l<ke.length;l+=1){let s=ht(d,ke,l),y=dt(s);ze.set(y,U[l]=vt(y,s))}ae=new Tt({props:{content:"?expand=relField1,relField2.subRelField21"}}),se=new Pt({});let ge=ie(d[3]);const ot=l=>l[10].code;for(let l=0;l<ge.length;l+=1){let s=_t(d,ge,l),y=ot(s);at.set(y,V[l]=wt(y,s))}let _e=ie(d[3]);const rt=l=>l[10].code;for(let l=0;l<_e.length;l+=1){let s=mt(d,_e,l),y=rt(s);it.set(y,B[l]=Ct(y,s))}return{c(){e=i("h3"),t=h("Update ("),m=h(a),p=h(")"),c=f(),u=i("div"),b=i("p"),O=h("Update a single "),T=i("strong"),H=h(M),S=h(" record."),E=f(),q=i("p"),q.innerHTML=`Body parameters could be sent as <code>application/json</code> or
        <code>multipart/form-data</code>.`,D=f(),J=i("p"),J.innerHTML=`File upload is supported only via <code>multipart/form-data</code>.
        <br/>
        For more info and examples you could check the detailed
        <a href="https://pocketbase.io/docs/files-handling" target="_blank" rel="noopener noreferrer">Files upload and handling docs
        </a>.`,I=f(),N&&N.c(),$=f(),ve(R.$$.fragment),L=f(),g=i("h6"),g.textContent="API details",v=f(),w=i("div"),z=i("strong"),z.textContent="PATCH",Q=f(),F=i("div"),_=i("p"),C=h("/api/collections/"),le=i("strong"),ne=h(W),Ue=h("/records/"),Oe=i("strong"),Oe.textContent=":id",Ve=f(),A&&A.c(),Se=f(),oe=i("div"),oe.textContent="Path parameters",Me=f(),re=i("table"),re.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr></thead> <tbody><tr><td>id</td> <td><span class="label">String</span></td> <td>ID of the record to update.</td></tr></tbody>',$e=f(),ce=i("div"),ce.textContent="Body Parameters",qe=f(),Y=i("table"),Re=i("thead"),Re.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr>',xe=f(),G=i("tbody"),P&&P.c(),He=f();for(let l=0;l<U.length;l+=1)U[l].c();De=f(),ue=i("div"),ue.textContent="Query parameters",Le=f(),K=i("table"),Pe=i("thead"),Pe.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr>',Qe=f(),pe=i("tbody"),X=i("tr"),Fe=i("td"),Fe.textContent="expand",We=f(),je=i("td"),je.innerHTML='<span class="label">String</span>',Ye=f(),j=i("td"),Ge=h(`Auto expand relations when returning the updated record. Ex.:
                `),ve(ae.$$.fragment),Ke=h(`
                Supports up to 6-levels depth nested relations expansion. `),Xe=i("br"),Ze=h(`
                The expanded relations will be appended to the record under the
                `),Be=i("code"),Be.textContent="expand",et=h(" property (eg. "),Ne=i("code"),Ne.textContent='"expand": {"relField1": {...}, ...}',tt=h(`). Only
                the relations that the user has permissions to `),Ae=i("strong"),Ae.textContent="view",lt=h(" will be expanded."),nt=f(),ve(se.$$.fragment),Ee=f(),fe=i("div"),fe.textContent="Responses",Ie=f(),Z=i("div"),be=i("div");for(let l=0;l<V.length;l+=1)V[l].c();st=f(),me=i("div");for(let l=0;l<B.length;l+=1)B[l].c();k(e,"class","m-b-sm"),k(u,"class","content txt-lg m-b-sm"),k(g,"class","m-b-xs"),k(z,"class","label label-primary"),k(F,"class","content"),k(w,"class","alert alert-warning"),k(oe,"class","section-title"),k(re,"class","table-compact table-border m-b-base"),k(ce,"class","section-title"),k(Y,"class","table-compact table-border m-b-base"),k(ue,"class","section-title"),k(K,"class","table-compact table-border m-b-lg"),k(fe,"class","section-title"),k(be,"class","tabs-header compact combined left"),k(me,"class","tabs-content"),k(Z,"class","tabs")},m(l,s){o(l,e,s),n(e,t),n(e,m),n(e,p),o(l,c,s),o(l,u,s),n(u,b),n(b,O),n(b,T),n(T,H),n(b,S),n(u,E),n(u,q),n(u,D),n(u,J),n(u,I),N&&N.m(u,null),o(l,$,s),we(R,l,s),o(l,L,s),o(l,g,s),o(l,v,s),o(l,w,s),n(w,z),n(w,Q),n(w,F),n(F,_),n(_,C),n(_,le),n(le,ne),n(_,Ue),n(_,Oe),n(w,Ve),A&&A.m(w,null),o(l,Se,s),o(l,oe,s),o(l,Me,s),o(l,re,s),o(l,$e,s),o(l,ce,s),o(l,qe,s),o(l,Y,s),n(Y,Re),n(Y,xe),n(Y,G),P&&P.m(G,null),n(G,He);for(let y=0;y<U.length;y+=1)U[y]&&U[y].m(G,null);o(l,De,s),o(l,ue,s),o(l,Le,s),o(l,K,s),n(K,Pe),n(K,Qe),n(K,pe),n(pe,X),n(X,Fe),n(X,We),n(X,je),n(X,Ye),n(X,j),n(j,Ge),we(ae,j,null),n(j,Ke),n(j,Xe),n(j,Ze),n(j,Be),n(j,et),n(j,Ne),n(j,tt),n(j,Ae),n(j,lt),n(pe,nt),we(se,pe,null),o(l,Ee,s),o(l,fe,s),o(l,Ie,s),o(l,Z,s),n(Z,be);for(let y=0;y<V.length;y+=1)V[y]&&V[y].m(be,null);n(Z,st),n(Z,me);for(let y=0;y<B.length;y+=1)B[y]&&B[y].m(me,null);ee=!0},p(l,[s]){var pt,ft;(!ee||s&1)&&a!==(a=l[0].name+"")&&te(m,a),(!ee||s&1)&&M!==(M=l[0].name+"")&&te(H,M),l[1]?N||(N=yt(),N.c(),N.m(u,null)):N&&(N.d(1),N=null);const y={};s&49&&(y.js=`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${l[5]}');

...

// example update data
const data = ${JSON.stringify(Object.assign({},l[4],x.dummyCollectionSchemaData(l[0],!0)),null,4)};

const record = await pb.collection('${(pt=l[0])==null?void 0:pt.name}').update('RECORD_ID', data);
    `),s&49&&(y.dart=`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${l[5]}');

...

// example update body
final body = <String, dynamic>${JSON.stringify(Object.assign({},l[4],x.dummyCollectionSchemaData(l[0],!0)),null,2)};

final record = await pb.collection('${(ft=l[0])==null?void 0:ft.name}').update('RECORD_ID', body: body);
    `),R.$set(y),(!ee||s&1)&&W!==(W=l[0].name+"")&&te(ne,W),l[7]?A||(A=kt(),A.c(),A.m(w,null)):A&&(A.d(1),A=null),l[1]?P?P.p(l,s):(P=gt(l),P.c(),P.m(G,He)):P&&(P.d(1),P=null),s&64&&(ke=ie(l[6]),U=Je(U,s,dt,1,l,ke,ze,G,bt,vt,null,ht)),s&12&&(ge=ie(l[3]),V=Je(V,s,ot,1,l,ge,at,be,bt,wt,null,_t)),s&12&&(_e=ie(l[3]),qt(),B=Je(B,s,rt,1,l,_e,it,me,Rt,Ct,null,mt),Ht())},i(l){if(!ee){he(R.$$.fragment,l),he(ae.$$.fragment,l),he(se.$$.fragment,l);for(let s=0;s<_e.length;s+=1)he(B[s]);ee=!0}},o(l){ye(R.$$.fragment,l),ye(ae.$$.fragment,l),ye(se.$$.fragment,l);for(let s=0;s<B.length;s+=1)ye(B[s]);ee=!1},d(l){l&&(r(e),r(c),r(u),r($),r(L),r(g),r(v),r(w),r(Se),r(oe),r(Me),r(re),r($e),r(ce),r(qe),r(Y),r(De),r(ue),r(Le),r(K),r(Ee),r(fe),r(Ie),r(Z)),N&&N.d(),Ce(R,l),A&&A.d(),P&&P.d();for(let s=0;s<U.length;s+=1)U[s].d();Ce(ae),Ce(se);for(let s=0;s<V.length;s+=1)V[s].d();for(let s=0;s<B.length;s+=1)B[s].d()}}}const Qt=d=>d.name=="emailVisibility";function Wt(d,e,t){let a,m,p,c,u,{collection:b}=e,O=200,T=[],M={};const H=S=>t(2,O=S.code);return d.$$set=S=>{"collection"in S&&t(0,b=S.collection)},d.$$.update=()=>{var S,E,q;d.$$.dirty&1&&t(1,a=(b==null?void 0:b.type)==="auth"),d.$$.dirty&1&&t(7,m=(b==null?void 0:b.updateRule)===null),d.$$.dirty&2&&t(8,p=a?["id","password","verified","email","emailVisibility"]:["id"]),d.$$.dirty&257&&t(6,c=((S=b==null?void 0:b.fields)==null?void 0:S.filter(D=>!D.hidden&&D.type!="autodate"&&!p.includes(D.name)))||[]),d.$$.dirty&1&&t(3,T=[{code:200,body:JSON.stringify(x.dummyCollectionRecord(b),null,2)},{code:400,body:`
                {
                  "code": 400,
                  "message": "Failed to update record.",
                  "data": {
                    "${(q=(E=b==null?void 0:b.fields)==null?void 0:E[0])==null?void 0:q.name}": {
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
            `},{code:404,body:`
                {
                  "code": 404,
                  "message": "The requested resource wasn't found.",
                  "data": {}
                }
            `}]),d.$$.dirty&2&&(a?t(4,M={password:"87654321",passwordConfirm:"87654321",oldPassword:"12345678"}):t(4,M={}))},t(5,u=x.getApiExampleUrl(Dt.baseURL)),[b,a,O,T,M,u,c,m,p,H]}class Kt extends Ot{constructor(e){super(),St(this,e,Wt,zt,Mt,{collection:0})}}export{Kt as default};
