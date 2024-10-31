import{S as Ze,i as tl,s as el,e,b as s,E as sl,f as a,g as m,u as ll,y as Ue,o as h,w as g,h as t,Q as nl,R as ve,T as se,c as Ut,m as jt,x as ke,U as je,V as ol,k as al,W as il,n as rl,t as $t,a as Ct,d as Qt,X as cl,C as Le,p as dl,r as Fe}from"./index-D5-_lkfc.js";import{F as pl}from"./FieldsQueryParam-D8rkBfEn.js";function fl(r){let n,o,i;return{c(){n=e("span"),n.textContent="Show details",o=s(),i=e("i"),a(n,"class","txt"),a(i,"class","ri-arrow-down-s-line")},m(f,b){m(f,n,b),m(f,o,b),m(f,i,b)},d(f){f&&(h(n),h(o),h(i))}}}function ul(r){let n,o,i;return{c(){n=e("span"),n.textContent="Hide details",o=s(),i=e("i"),a(n,"class","txt"),a(i,"class","ri-arrow-up-s-line")},m(f,b){m(f,n,b),m(f,o,b),m(f,i,b)},d(f){f&&(h(n),h(o),h(i))}}}function Qe(r){let n,o,i,f,b,p,u,C,_,$,d,tt,kt,zt,S,Kt,H,rt,R,et,ne,U,j,oe,ct,yt,lt,vt,ae,dt,pt,st,N,Jt,Ft,y,nt,Lt,Vt,At,Q,ot,Tt,Wt,Pt,F,ft,Rt,ie,ut,re,M,Ot,at,Et,O,mt,ce,z,St,Xt,Nt,de,q,Yt,K,ht,pe,I,fe,B,ue,P,qt,J,bt,me,gt,he,x,Dt,it,Ht,be,Mt,Zt,V,_t,ge,It,_e,wt,we,W,G,xe,xt,te,X,ee,L,Y,E,Bt,$e,Z,v,Gt;return{c(){n=e("p"),n.innerHTML=`The syntax basically follows the format
        <code><span class="txt-success">OPERAND</span> <span class="txt-danger">OPERATOR</span> <span class="txt-success">OPERAND</span></code>, where:`,o=s(),i=e("ul"),f=e("li"),f.innerHTML=`<code class="txt-success">OPERAND</code> - could be any of the above field literal, string (single
            or double quoted), number, null, true, false`,b=s(),p=e("li"),u=e("code"),u.textContent="OPERATOR",C=g(` - is one of:
            `),_=e("br"),$=s(),d=e("ul"),tt=e("li"),kt=e("code"),kt.textContent="=",zt=s(),S=e("span"),S.textContent="Equal",Kt=s(),H=e("li"),rt=e("code"),rt.textContent="!=",R=s(),et=e("span"),et.textContent="NOT equal",ne=s(),U=e("li"),j=e("code"),j.textContent=">",oe=s(),ct=e("span"),ct.textContent="Greater than",yt=s(),lt=e("li"),vt=e("code"),vt.textContent=">=",ae=s(),dt=e("span"),dt.textContent="Greater than or equal",pt=s(),st=e("li"),N=e("code"),N.textContent="<",Jt=s(),Ft=e("span"),Ft.textContent="Less than",y=s(),nt=e("li"),Lt=e("code"),Lt.textContent="<=",Vt=s(),At=e("span"),At.textContent="Less than or equal",Q=s(),ot=e("li"),Tt=e("code"),Tt.textContent="~",Wt=s(),Pt=e("span"),Pt.textContent=`Like/Contains (if not specified auto wraps the right string OPERAND in a "%" for
                        wildcard match)`,F=s(),ft=e("li"),Rt=e("code"),Rt.textContent="!~",ie=s(),ut=e("span"),ut.textContent=`NOT Like/Contains (if not specified auto wraps the right string OPERAND in a "%" for
                        wildcard match)`,re=s(),M=e("li"),Ot=e("code"),Ot.textContent="?=",at=s(),Et=e("em"),Et.textContent="Any/At least one of",O=s(),mt=e("span"),mt.textContent="Equal",ce=s(),z=e("li"),St=e("code"),St.textContent="?!=",Xt=s(),Nt=e("em"),Nt.textContent="Any/At least one of",de=s(),q=e("span"),q.textContent="NOT equal",Yt=s(),K=e("li"),ht=e("code"),ht.textContent="?>",pe=s(),I=e("em"),I.textContent="Any/At least one of",fe=s(),B=e("span"),B.textContent="Greater than",ue=s(),P=e("li"),qt=e("code"),qt.textContent="?>=",J=s(),bt=e("em"),bt.textContent="Any/At least one of",me=s(),gt=e("span"),gt.textContent="Greater than or equal",he=s(),x=e("li"),Dt=e("code"),Dt.textContent="?<",it=s(),Ht=e("em"),Ht.textContent="Any/At least one of",be=s(),Mt=e("span"),Mt.textContent="Less than",Zt=s(),V=e("li"),_t=e("code"),_t.textContent="?<=",ge=s(),It=e("em"),It.textContent="Any/At least one of",_e=s(),wt=e("span"),wt.textContent="Less than or equal",we=s(),W=e("li"),G=e("code"),G.textContent="?~",xe=s(),xt=e("em"),xt.textContent="Any/At least one of",te=s(),X=e("span"),X.textContent=`Like/Contains (if not specified auto wraps the right string OPERAND in a "%" for
                        wildcard match)`,ee=s(),L=e("li"),Y=e("code"),Y.textContent="?!~",E=s(),Bt=e("em"),Bt.textContent="Any/At least one of",$e=s(),Z=e("span"),Z.textContent=`NOT Like/Contains (if not specified auto wraps the right string OPERAND in a "%" for
                        wildcard match)`,v=s(),Gt=e("p"),Gt.innerHTML=`To group and combine several expressions you could use brackets
        <code>(...)</code>, <code>&amp;&amp;</code> (AND) and <code>||</code> (OR) tokens.`,a(u,"class","txt-danger"),a(kt,"class","filter-op svelte-1w7s5nw"),a(S,"class","txt"),a(rt,"class","filter-op svelte-1w7s5nw"),a(et,"class","txt"),a(j,"class","filter-op svelte-1w7s5nw"),a(ct,"class","txt"),a(vt,"class","filter-op svelte-1w7s5nw"),a(dt,"class","txt"),a(N,"class","filter-op svelte-1w7s5nw"),a(Ft,"class","txt"),a(Lt,"class","filter-op svelte-1w7s5nw"),a(At,"class","txt"),a(Tt,"class","filter-op svelte-1w7s5nw"),a(Pt,"class","txt"),a(Rt,"class","filter-op svelte-1w7s5nw"),a(ut,"class","txt"),a(Ot,"class","filter-op svelte-1w7s5nw"),a(Et,"class","txt-hint"),a(mt,"class","txt"),a(St,"class","filter-op svelte-1w7s5nw"),a(Nt,"class","txt-hint"),a(q,"class","txt"),a(ht,"class","filter-op svelte-1w7s5nw"),a(I,"class","txt-hint"),a(B,"class","txt"),a(qt,"class","filter-op svelte-1w7s5nw"),a(bt,"class","txt-hint"),a(gt,"class","txt"),a(Dt,"class","filter-op svelte-1w7s5nw"),a(Ht,"class","txt-hint"),a(Mt,"class","txt"),a(_t,"class","filter-op svelte-1w7s5nw"),a(It,"class","txt-hint"),a(wt,"class","txt"),a(G,"class","filter-op svelte-1w7s5nw"),a(xt,"class","txt-hint"),a(X,"class","txt"),a(Y,"class","filter-op svelte-1w7s5nw"),a(Bt,"class","txt-hint"),a(Z,"class","txt")},m(A,k){m(A,n,k),m(A,o,k),m(A,i,k),t(i,f),t(i,b),t(i,p),t(p,u),t(p,C),t(p,_),t(p,$),t(p,d),t(d,tt),t(tt,kt),t(tt,zt),t(tt,S),t(d,Kt),t(d,H),t(H,rt),t(H,R),t(H,et),t(d,ne),t(d,U),t(U,j),t(U,oe),t(U,ct),t(d,yt),t(d,lt),t(lt,vt),t(lt,ae),t(lt,dt),t(d,pt),t(d,st),t(st,N),t(st,Jt),t(st,Ft),t(d,y),t(d,nt),t(nt,Lt),t(nt,Vt),t(nt,At),t(d,Q),t(d,ot),t(ot,Tt),t(ot,Wt),t(ot,Pt),t(d,F),t(d,ft),t(ft,Rt),t(ft,ie),t(ft,ut),t(d,re),t(d,M),t(M,Ot),t(M,at),t(M,Et),t(M,O),t(M,mt),t(d,ce),t(d,z),t(z,St),t(z,Xt),t(z,Nt),t(z,de),t(z,q),t(d,Yt),t(d,K),t(K,ht),t(K,pe),t(K,I),t(K,fe),t(K,B),t(d,ue),t(d,P),t(P,qt),t(P,J),t(P,bt),t(P,me),t(P,gt),t(d,he),t(d,x),t(x,Dt),t(x,it),t(x,Ht),t(x,be),t(x,Mt),t(d,Zt),t(d,V),t(V,_t),t(V,ge),t(V,It),t(V,_e),t(V,wt),t(d,we),t(d,W),t(W,G),t(W,xe),t(W,xt),t(W,te),t(W,X),t(d,ee),t(d,L),t(L,Y),t(L,E),t(L,Bt),t(L,$e),t(L,Z),m(A,v,k),m(A,Gt,k)},d(A){A&&(h(n),h(o),h(i),h(v),h(Gt))}}}function ml(r){let n,o,i,f,b;function p($,d){return $[0]?ul:fl}let u=p(r),C=u(r),_=r[0]&&Qe();return{c(){n=e("button"),C.c(),o=s(),_&&_.c(),i=sl(),a(n,"class","btn btn-sm btn-secondary m-t-10")},m($,d){m($,n,d),C.m(n,null),m($,o,d),_&&_.m($,d),m($,i,d),f||(b=ll(n,"click",r[1]),f=!0)},p($,[d]){u!==(u=p($))&&(C.d(1),C=u($),C&&(C.c(),C.m(n,null))),$[0]?_||(_=Qe(),_.c(),_.m(i.parentNode,i)):_&&(_.d(1),_=null)},i:Ue,o:Ue,d($){$&&(h(n),h(o),h(i)),C.d(),_&&_.d($),f=!1,b()}}}function hl(r,n,o){let i=!1;function f(){o(0,i=!i)}return[i,f]}class bl extends Ze{constructor(n){super(),tl(this,n,hl,ml,el,{})}}function ze(r,n,o){const i=r.slice();return i[8]=n[o],i}function Ke(r,n,o){const i=r.slice();return i[8]=n[o],i}function Je(r,n,o){const i=r.slice();return i[13]=n[o],i[15]=o,i}function Ve(r){let n;return{c(){n=e("p"),n.innerHTML="Requires superuser <code>Authorization:TOKEN</code> header",a(n,"class","txt-hint txt-sm txt-right")},m(o,i){m(o,n,i)},d(o){o&&h(n)}}}function We(r){let n,o=r[13]+"",i,f=r[15]<r[4].length-1?", ":"",b;return{c(){n=e("code"),i=g(o),b=g(f)},m(p,u){m(p,n,u),t(n,i),m(p,b,u)},p(p,u){u&16&&o!==(o=p[13]+"")&&ke(i,o),u&16&&f!==(f=p[15]<p[4].length-1?", ":"")&&ke(b,f)},d(p){p&&(h(n),h(b))}}}function Xe(r,n){let o,i,f;function b(){return n[7](n[8])}return{key:r,first:null,c(){o=e("button"),o.textContent=`${n[8].code} `,a(o,"type","button"),a(o,"class","tab-item"),Fe(o,"active",n[2]===n[8].code),this.first=o},m(p,u){m(p,o,u),i||(f=ll(o,"click",b),i=!0)},p(p,u){n=p,u&36&&Fe(o,"active",n[2]===n[8].code)},d(p){p&&h(o),i=!1,f()}}}function Ye(r,n){let o,i,f,b;return i=new ve({props:{content:n[8].body}}),{key:r,first:null,c(){o=e("div"),Ut(i.$$.fragment),f=s(),a(o,"class","tab-item"),Fe(o,"active",n[2]===n[8].code),this.first=o},m(p,u){m(p,o,u),jt(i,o,null),t(o,f),b=!0},p(p,u){n=p,(!b||u&36)&&Fe(o,"active",n[2]===n[8].code)},i(p){b||($t(i.$$.fragment,p),b=!0)},o(p){Ct(i.$$.fragment,p),b=!1},d(p){p&&h(o),Qt(i)}}}function gl(r){var Pe,Re,Oe,Ee,Se,Ne;let n,o,i=r[0].name+"",f,b,p,u,C,_,$,d=r[0].name+"",tt,kt,zt,S,Kt,H,rt,R,et,ne,U,j,oe,ct,yt=r[0].name+"",lt,vt,ae,dt,pt,st,N,Jt,Ft,y,nt,Lt,Vt,At,Q,ot,Tt,Wt,Pt,F,ft,Rt,ie,ut,re,M,Ot,at,Et,O,mt,ce,z,St,Xt,Nt,de,q,Yt,K,ht,pe,I,fe,B,ue,P,qt,J,bt,me,gt,he,x,Dt,it,Ht,be,Mt,Zt,V,_t,ge,It,_e,wt,we,W,G,xe,xt,te,X,ee,L,Y,E=[],Bt=new Map,$e,Z,v=[],Gt=new Map,A;S=new nl({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${r[3]}');

        ...

        // fetch a paginated records list
        const resultList = await pb.collection('${(Pe=r[0])==null?void 0:Pe.name}').getList(1, 50, {
            filter: 'created >= "2022-01-01 00:00:00" && someField1 != someField2',
        });

        // you can also fetch all records at once via getFullList
        const records = await pb.collection('${(Re=r[0])==null?void 0:Re.name}').getFullList({
            sort: '-created',
        });

        // or fetch only the first record that matches the specified filter
        const record = await pb.collection('${(Oe=r[0])==null?void 0:Oe.name}').getFirstListItem('someField="test"', {
            expand: 'relField1,relField2.subRelField',
        });
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${r[3]}');

        ...

        // fetch a paginated records list
        final resultList = await pb.collection('${(Ee=r[0])==null?void 0:Ee.name}').getList(
          page: 1,
          perPage: 50,
          filter: 'created >= "2022-01-01 00:00:00" && someField1 != someField2',
        );

        // you can also fetch all records at once via getFullList
        final records = await pb.collection('${(Se=r[0])==null?void 0:Se.name}').getFullList(
          sort: '-created',
        );

        // or fetch only the first record that matches the specified filter
        final record = await pb.collection('${(Ne=r[0])==null?void 0:Ne.name}').getFirstListItem(
          'someField="test"',
          expand: 'relField1,relField2.subRelField',
        );
    `}});let k=r[1]&&Ve();at=new ve({props:{content:`
                        // DESC by created and ASC by id
                        ?sort=-created,id
                    `}});let le=se(r[4]),T=[];for(let l=0;l<le.length;l+=1)T[l]=We(Je(r,le,l));B=new ve({props:{content:`
                        ?filter=(id='abc' && created>'2022-01-01')
                    `}}),P=new bl({}),it=new ve({props:{content:"?expand=relField1,relField2.subRelField"}}),G=new pl({});let ye=se(r[5]);const Ae=l=>l[8].code;for(let l=0;l<ye.length;l+=1){let c=Ke(r,ye,l),w=Ae(c);Bt.set(w,E[l]=Xe(w,c))}let Ce=se(r[5]);const Te=l=>l[8].code;for(let l=0;l<Ce.length;l+=1){let c=ze(r,Ce,l),w=Te(c);Gt.set(w,v[l]=Ye(w,c))}return{c(){n=e("h3"),o=g("List/Search ("),f=g(i),b=g(")"),p=s(),u=e("div"),C=e("p"),_=g("Fetch a paginated "),$=e("strong"),tt=g(d),kt=g(" records list, supporting sorting and filtering."),zt=s(),Ut(S.$$.fragment),Kt=s(),H=e("h6"),H.textContent="API details",rt=s(),R=e("div"),et=e("strong"),et.textContent="GET",ne=s(),U=e("div"),j=e("p"),oe=g("/api/collections/"),ct=e("strong"),lt=g(yt),vt=g("/records"),ae=s(),k&&k.c(),dt=s(),pt=e("div"),pt.textContent="Query parameters",st=s(),N=e("table"),Jt=e("thead"),Jt.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr>',Ft=s(),y=e("tbody"),nt=e("tr"),nt.innerHTML='<td>page</td> <td><span class="label">Number</span></td> <td>The page (aka. offset) of the paginated list (default to 1).</td>',Lt=s(),Vt=e("tr"),Vt.innerHTML='<td>perPage</td> <td><span class="label">Number</span></td> <td>Specify the max returned records per page (default to 30).</td>',At=s(),Q=e("tr"),ot=e("td"),ot.textContent="sort",Tt=s(),Wt=e("td"),Wt.innerHTML='<span class="label">String</span>',Pt=s(),F=e("td"),ft=g("Specify the records order attribute(s). "),Rt=e("br"),ie=g(`
                Add `),ut=e("code"),ut.textContent="-",re=g(" / "),M=e("code"),M.textContent="+",Ot=g(` (default) in front of the attribute for DESC / ASC order.
                Ex.:
                `),Ut(at.$$.fragment),Et=s(),O=e("p"),mt=e("strong"),mt.textContent="Supported record sort fields:",ce=s(),z=e("br"),St=s(),Xt=e("code"),Xt.textContent="@random",Nt=g(`,
                    `);for(let l=0;l<T.length;l+=1)T[l].c();de=s(),q=e("tr"),Yt=e("td"),Yt.textContent="filter",K=s(),ht=e("td"),ht.innerHTML='<span class="label">String</span>',pe=s(),I=e("td"),fe=g(`Filter the returned records. Ex.:
                `),Ut(B.$$.fragment),ue=s(),Ut(P.$$.fragment),qt=s(),J=e("tr"),bt=e("td"),bt.textContent="expand",me=s(),gt=e("td"),gt.innerHTML='<span class="label">String</span>',he=s(),x=e("td"),Dt=g(`Auto expand record relations. Ex.:
                `),Ut(it.$$.fragment),Ht=g(`
                Supports up to 6-levels depth nested relations expansion. `),be=e("br"),Mt=g(`
                The expanded relations will be appended to each individual record under the
                `),Zt=e("code"),Zt.textContent="expand",V=g(" property (eg. "),_t=e("code"),_t.textContent='"expand": {"relField1": {...}, ...}',ge=g(`).
                `),It=e("br"),_e=g(`
                Only the relations to which the request user has permissions to `),wt=e("strong"),wt.textContent="view",we=g(" will be expanded."),W=s(),Ut(G.$$.fragment),xe=s(),xt=e("tr"),xt.innerHTML=`<td id="query-page">skipTotal</td> <td><span class="label">Boolean</span></td> <td>If it is set the total counts query will be skipped and the response fields
                <code>totalItems</code> and <code>totalPages</code> will have <code>-1</code> value.
                <br/>
                This could drastically speed up the search queries when the total counters are not needed or cursor
                based pagination is used.
                <br/>
                For optimization purposes, it is set by default for the
                <code>getFirstListItem()</code>
                and
                <code>getFullList()</code> SDKs methods.</td>`,te=s(),X=e("div"),X.textContent="Responses",ee=s(),L=e("div"),Y=e("div");for(let l=0;l<E.length;l+=1)E[l].c();$e=s(),Z=e("div");for(let l=0;l<v.length;l+=1)v[l].c();a(n,"class","m-b-sm"),a(u,"class","content txt-lg m-b-sm"),a(H,"class","m-b-xs"),a(et,"class","label label-primary"),a(U,"class","content"),a(R,"class","alert alert-info"),a(pt,"class","section-title"),a(N,"class","table-compact table-border m-b-base"),a(X,"class","section-title"),a(Y,"class","tabs-header compact combined left"),a(Z,"class","tabs-content"),a(L,"class","tabs")},m(l,c){m(l,n,c),t(n,o),t(n,f),t(n,b),m(l,p,c),m(l,u,c),t(u,C),t(C,_),t(C,$),t($,tt),t(C,kt),m(l,zt,c),jt(S,l,c),m(l,Kt,c),m(l,H,c),m(l,rt,c),m(l,R,c),t(R,et),t(R,ne),t(R,U),t(U,j),t(j,oe),t(j,ct),t(ct,lt),t(j,vt),t(R,ae),k&&k.m(R,null),m(l,dt,c),m(l,pt,c),m(l,st,c),m(l,N,c),t(N,Jt),t(N,Ft),t(N,y),t(y,nt),t(y,Lt),t(y,Vt),t(y,At),t(y,Q),t(Q,ot),t(Q,Tt),t(Q,Wt),t(Q,Pt),t(Q,F),t(F,ft),t(F,Rt),t(F,ie),t(F,ut),t(F,re),t(F,M),t(F,Ot),jt(at,F,null),t(F,Et),t(F,O),t(O,mt),t(O,ce),t(O,z),t(O,St),t(O,Xt),t(O,Nt);for(let w=0;w<T.length;w+=1)T[w]&&T[w].m(O,null);t(y,de),t(y,q),t(q,Yt),t(q,K),t(q,ht),t(q,pe),t(q,I),t(I,fe),jt(B,I,null),t(I,ue),jt(P,I,null),t(y,qt),t(y,J),t(J,bt),t(J,me),t(J,gt),t(J,he),t(J,x),t(x,Dt),jt(it,x,null),t(x,Ht),t(x,be),t(x,Mt),t(x,Zt),t(x,V),t(x,_t),t(x,ge),t(x,It),t(x,_e),t(x,wt),t(x,we),t(y,W),jt(G,y,null),t(y,xe),t(y,xt),m(l,te,c),m(l,X,c),m(l,ee,c),m(l,L,c),t(L,Y);for(let w=0;w<E.length;w+=1)E[w]&&E[w].m(Y,null);t(L,$e),t(L,Z);for(let w=0;w<v.length;w+=1)v[w]&&v[w].m(Z,null);A=!0},p(l,[c]){var qe,De,He,Me,Ie,Be;(!A||c&1)&&i!==(i=l[0].name+"")&&ke(f,i),(!A||c&1)&&d!==(d=l[0].name+"")&&ke(tt,d);const w={};if(c&9&&(w.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${l[3]}');

        ...

        // fetch a paginated records list
        const resultList = await pb.collection('${(qe=l[0])==null?void 0:qe.name}').getList(1, 50, {
            filter: 'created >= "2022-01-01 00:00:00" && someField1 != someField2',
        });

        // you can also fetch all records at once via getFullList
        const records = await pb.collection('${(De=l[0])==null?void 0:De.name}').getFullList({
            sort: '-created',
        });

        // or fetch only the first record that matches the specified filter
        const record = await pb.collection('${(He=l[0])==null?void 0:He.name}').getFirstListItem('someField="test"', {
            expand: 'relField1,relField2.subRelField',
        });
    `),c&9&&(w.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${l[3]}');

        ...

        // fetch a paginated records list
        final resultList = await pb.collection('${(Me=l[0])==null?void 0:Me.name}').getList(
          page: 1,
          perPage: 50,
          filter: 'created >= "2022-01-01 00:00:00" && someField1 != someField2',
        );

        // you can also fetch all records at once via getFullList
        final records = await pb.collection('${(Ie=l[0])==null?void 0:Ie.name}').getFullList(
          sort: '-created',
        );

        // or fetch only the first record that matches the specified filter
        final record = await pb.collection('${(Be=l[0])==null?void 0:Be.name}').getFirstListItem(
          'someField="test"',
          expand: 'relField1,relField2.subRelField',
        );
    `),S.$set(w),(!A||c&1)&&yt!==(yt=l[0].name+"")&&ke(lt,yt),l[1]?k||(k=Ve(),k.c(),k.m(R,null)):k&&(k.d(1),k=null),c&16){le=se(l[4]);let D;for(D=0;D<le.length;D+=1){const Ge=Je(l,le,D);T[D]?T[D].p(Ge,c):(T[D]=We(Ge),T[D].c(),T[D].m(O,null))}for(;D<T.length;D+=1)T[D].d(1);T.length=le.length}c&36&&(ye=se(l[5]),E=je(E,c,Ae,1,l,ye,Bt,Y,ol,Xe,null,Ke)),c&36&&(Ce=se(l[5]),al(),v=je(v,c,Te,1,l,Ce,Gt,Z,il,Ye,null,ze),rl())},i(l){if(!A){$t(S.$$.fragment,l),$t(at.$$.fragment,l),$t(B.$$.fragment,l),$t(P.$$.fragment,l),$t(it.$$.fragment,l),$t(G.$$.fragment,l);for(let c=0;c<Ce.length;c+=1)$t(v[c]);A=!0}},o(l){Ct(S.$$.fragment,l),Ct(at.$$.fragment,l),Ct(B.$$.fragment,l),Ct(P.$$.fragment,l),Ct(it.$$.fragment,l),Ct(G.$$.fragment,l);for(let c=0;c<v.length;c+=1)Ct(v[c]);A=!1},d(l){l&&(h(n),h(p),h(u),h(zt),h(Kt),h(H),h(rt),h(R),h(dt),h(pt),h(st),h(N),h(te),h(X),h(ee),h(L)),Qt(S,l),k&&k.d(),Qt(at),cl(T,l),Qt(B),Qt(P),Qt(it),Qt(G);for(let c=0;c<E.length;c+=1)E[c].d();for(let c=0;c<v.length;c+=1)v[c].d()}}}function _l(r,n,o){let i,f,b,p,{collection:u}=n,C=200,_=[];const $=d=>o(2,C=d.code);return r.$$set=d=>{"collection"in d&&o(0,u=d.collection)},r.$$.update=()=>{r.$$.dirty&1&&o(4,i=Le.getAllCollectionIdentifiers(u)),r.$$.dirty&1&&o(1,f=(u==null?void 0:u.listRule)===null),r.$$.dirty&1&&o(6,p=Le.dummyCollectionRecord(u)),r.$$.dirty&67&&u!=null&&u.id&&(_.push({code:200,body:JSON.stringify({page:1,perPage:30,totalPages:1,totalItems:2,items:[p,Object.assign({},p,{id:p+"2"})]},null,2)}),_.push({code:400,body:`
                {
                  "code": 400,
                  "message": "Something went wrong while processing your request. Invalid filter.",
                  "data": {}
                }
            `}),f&&_.push({code:403,body:`
                    {
                      "code": 403,
                      "message": "Only superusers can access this action.",
                      "data": {}
                    }
                `}))},o(3,b=Le.getApiExampleUrl(dl.baseURL)),[u,f,C,b,i,_,p,$]}class $l extends Ze{constructor(n){super(),tl(this,n,_l,gl,el,{collection:0})}}export{$l as default};
