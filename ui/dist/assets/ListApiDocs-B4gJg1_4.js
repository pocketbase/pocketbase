import{S as Ze,i as tl,s as el,e,b as s,E as sl,f as a,g as u,r as ll,x as Qe,o as m,v as _,h as t,N as Fe,O as se,c as Qt,m as Ut,w as ke,P as Ue,Q as nl,k as ol,R as al,n as il,t as $t,a as Ct,d as jt,T as rl,C as ve,A as cl,q as Le}from"./index-D0DO79Dq.js";import{S as dl}from"./SdkTabs-DC6EUYpr.js";import{F as pl}from"./FieldsQueryParam-BwleQAus.js";function fl(d){let n,o,i;return{c(){n=e("span"),n.textContent="Show details",o=s(),i=e("i"),a(n,"class","txt"),a(i,"class","ri-arrow-down-s-line")},m(f,h){u(f,n,h),u(f,o,h),u(f,i,h)},d(f){f&&(m(n),m(o),m(i))}}}function ul(d){let n,o,i;return{c(){n=e("span"),n.textContent="Hide details",o=s(),i=e("i"),a(n,"class","txt"),a(i,"class","ri-arrow-up-s-line")},m(f,h){u(f,n,h),u(f,o,h),u(f,i,h)},d(f){f&&(m(n),m(o),m(i))}}}function je(d){let n,o,i,f,h,r,b,$,C,g,p,tt,kt,zt,E,Kt,H,rt,R,et,ne,Q,U,oe,ct,yt,lt,vt,ae,dt,pt,st,N,Jt,Ft,y,nt,Lt,Vt,At,j,ot,Tt,Wt,Pt,F,ft,Rt,ie,ut,re,M,Ot,at,St,O,mt,ce,z,Et,Xt,Nt,de,q,Yt,K,ht,pe,I,fe,B,ue,P,qt,J,bt,me,gt,he,x,Dt,it,Ht,be,Mt,Zt,V,_t,ge,It,_e,wt,we,W,G,xe,xt,te,X,ee,L,Y,S,Bt,$e,Z,v,Gt;return{c(){n=e("p"),n.innerHTML=`The syntax basically follows the format
        <code><span class="txt-success">OPERAND</span> <span class="txt-danger">OPERATOR</span> <span class="txt-success">OPERAND</span></code>, where:`,o=s(),i=e("ul"),f=e("li"),f.innerHTML=`<code class="txt-success">OPERAND</code> - could be any of the above field literal, string (single
            or double quoted), number, null, true, false`,h=s(),r=e("li"),b=e("code"),b.textContent="OPERATOR",$=_(` - is one of:
            `),C=e("br"),g=s(),p=e("ul"),tt=e("li"),kt=e("code"),kt.textContent="=",zt=s(),E=e("span"),E.textContent="Equal",Kt=s(),H=e("li"),rt=e("code"),rt.textContent="!=",R=s(),et=e("span"),et.textContent="NOT equal",ne=s(),Q=e("li"),U=e("code"),U.textContent=">",oe=s(),ct=e("span"),ct.textContent="Greater than",yt=s(),lt=e("li"),vt=e("code"),vt.textContent=">=",ae=s(),dt=e("span"),dt.textContent="Greater than or equal",pt=s(),st=e("li"),N=e("code"),N.textContent="<",Jt=s(),Ft=e("span"),Ft.textContent="Less than",y=s(),nt=e("li"),Lt=e("code"),Lt.textContent="<=",Vt=s(),At=e("span"),At.textContent="Less than or equal",j=s(),ot=e("li"),Tt=e("code"),Tt.textContent="~",Wt=s(),Pt=e("span"),Pt.textContent=`Like/Contains (if not specified auto wraps the right string OPERAND in a "%" for
                        wildcard match)`,F=s(),ft=e("li"),Rt=e("code"),Rt.textContent="!~",ie=s(),ut=e("span"),ut.textContent=`NOT Like/Contains (if not specified auto wraps the right string OPERAND in a "%" for
                        wildcard match)`,re=s(),M=e("li"),Ot=e("code"),Ot.textContent="?=",at=s(),St=e("em"),St.textContent="Any/At least one of",O=s(),mt=e("span"),mt.textContent="Equal",ce=s(),z=e("li"),Et=e("code"),Et.textContent="?!=",Xt=s(),Nt=e("em"),Nt.textContent="Any/At least one of",de=s(),q=e("span"),q.textContent="NOT equal",Yt=s(),K=e("li"),ht=e("code"),ht.textContent="?>",pe=s(),I=e("em"),I.textContent="Any/At least one of",fe=s(),B=e("span"),B.textContent="Greater than",ue=s(),P=e("li"),qt=e("code"),qt.textContent="?>=",J=s(),bt=e("em"),bt.textContent="Any/At least one of",me=s(),gt=e("span"),gt.textContent="Greater than or equal",he=s(),x=e("li"),Dt=e("code"),Dt.textContent="?<",it=s(),Ht=e("em"),Ht.textContent="Any/At least one of",be=s(),Mt=e("span"),Mt.textContent="Less than",Zt=s(),V=e("li"),_t=e("code"),_t.textContent="?<=",ge=s(),It=e("em"),It.textContent="Any/At least one of",_e=s(),wt=e("span"),wt.textContent="Less than or equal",we=s(),W=e("li"),G=e("code"),G.textContent="?~",xe=s(),xt=e("em"),xt.textContent="Any/At least one of",te=s(),X=e("span"),X.textContent=`Like/Contains (if not specified auto wraps the right string OPERAND in a "%" for
                        wildcard match)`,ee=s(),L=e("li"),Y=e("code"),Y.textContent="?!~",S=s(),Bt=e("em"),Bt.textContent="Any/At least one of",$e=s(),Z=e("span"),Z.textContent=`NOT Like/Contains (if not specified auto wraps the right string OPERAND in a "%" for
                        wildcard match)`,v=s(),Gt=e("p"),Gt.innerHTML=`To group and combine several expressions you could use brackets
        <code>(...)</code>, <code>&amp;&amp;</code> (AND) and <code>||</code> (OR) tokens.`,a(b,"class","txt-danger"),a(kt,"class","filter-op svelte-1w7s5nw"),a(E,"class","txt"),a(rt,"class","filter-op svelte-1w7s5nw"),a(et,"class","txt"),a(U,"class","filter-op svelte-1w7s5nw"),a(ct,"class","txt"),a(vt,"class","filter-op svelte-1w7s5nw"),a(dt,"class","txt"),a(N,"class","filter-op svelte-1w7s5nw"),a(Ft,"class","txt"),a(Lt,"class","filter-op svelte-1w7s5nw"),a(At,"class","txt"),a(Tt,"class","filter-op svelte-1w7s5nw"),a(Pt,"class","txt"),a(Rt,"class","filter-op svelte-1w7s5nw"),a(ut,"class","txt"),a(Ot,"class","filter-op svelte-1w7s5nw"),a(St,"class","txt-hint"),a(mt,"class","txt"),a(Et,"class","filter-op svelte-1w7s5nw"),a(Nt,"class","txt-hint"),a(q,"class","txt"),a(ht,"class","filter-op svelte-1w7s5nw"),a(I,"class","txt-hint"),a(B,"class","txt"),a(qt,"class","filter-op svelte-1w7s5nw"),a(bt,"class","txt-hint"),a(gt,"class","txt"),a(Dt,"class","filter-op svelte-1w7s5nw"),a(Ht,"class","txt-hint"),a(Mt,"class","txt"),a(_t,"class","filter-op svelte-1w7s5nw"),a(It,"class","txt-hint"),a(wt,"class","txt"),a(G,"class","filter-op svelte-1w7s5nw"),a(xt,"class","txt-hint"),a(X,"class","txt"),a(Y,"class","filter-op svelte-1w7s5nw"),a(Bt,"class","txt-hint"),a(Z,"class","txt")},m(A,k){u(A,n,k),u(A,o,k),u(A,i,k),t(i,f),t(i,h),t(i,r),t(r,b),t(r,$),t(r,C),t(r,g),t(r,p),t(p,tt),t(tt,kt),t(tt,zt),t(tt,E),t(p,Kt),t(p,H),t(H,rt),t(H,R),t(H,et),t(p,ne),t(p,Q),t(Q,U),t(Q,oe),t(Q,ct),t(p,yt),t(p,lt),t(lt,vt),t(lt,ae),t(lt,dt),t(p,pt),t(p,st),t(st,N),t(st,Jt),t(st,Ft),t(p,y),t(p,nt),t(nt,Lt),t(nt,Vt),t(nt,At),t(p,j),t(p,ot),t(ot,Tt),t(ot,Wt),t(ot,Pt),t(p,F),t(p,ft),t(ft,Rt),t(ft,ie),t(ft,ut),t(p,re),t(p,M),t(M,Ot),t(M,at),t(M,St),t(M,O),t(M,mt),t(p,ce),t(p,z),t(z,Et),t(z,Xt),t(z,Nt),t(z,de),t(z,q),t(p,Yt),t(p,K),t(K,ht),t(K,pe),t(K,I),t(K,fe),t(K,B),t(p,ue),t(p,P),t(P,qt),t(P,J),t(P,bt),t(P,me),t(P,gt),t(p,he),t(p,x),t(x,Dt),t(x,it),t(x,Ht),t(x,be),t(x,Mt),t(p,Zt),t(p,V),t(V,_t),t(V,ge),t(V,It),t(V,_e),t(V,wt),t(p,we),t(p,W),t(W,G),t(W,xe),t(W,xt),t(W,te),t(W,X),t(p,ee),t(p,L),t(L,Y),t(L,S),t(L,Bt),t(L,$e),t(L,Z),u(A,v,k),u(A,Gt,k)},d(A){A&&(m(n),m(o),m(i),m(v),m(Gt))}}}function ml(d){let n,o,i,f,h;function r(g,p){return g[0]?ul:fl}let b=r(d),$=b(d),C=d[0]&&je();return{c(){n=e("button"),$.c(),o=s(),C&&C.c(),i=sl(),a(n,"class","btn btn-sm btn-secondary m-t-10")},m(g,p){u(g,n,p),$.m(n,null),u(g,o,p),C&&C.m(g,p),u(g,i,p),f||(h=ll(n,"click",d[1]),f=!0)},p(g,[p]){b!==(b=r(g))&&($.d(1),$=b(g),$&&($.c(),$.m(n,null))),g[0]?C||(C=je(),C.c(),C.m(i.parentNode,i)):C&&(C.d(1),C=null)},i:Qe,o:Qe,d(g){g&&(m(n),m(o),m(i)),$.d(),C&&C.d(g),f=!1,h()}}}function hl(d,n,o){let i=!1;function f(){o(0,i=!i)}return[i,f]}class bl extends Ze{constructor(n){super(),tl(this,n,hl,ml,el,{})}}function ze(d,n,o){const i=d.slice();return i[7]=n[o],i}function Ke(d,n,o){const i=d.slice();return i[7]=n[o],i}function Je(d,n,o){const i=d.slice();return i[12]=n[o],i[14]=o,i}function Ve(d){let n;return{c(){n=e("p"),n.innerHTML="Requires admin <code>Authorization:TOKEN</code> header",a(n,"class","txt-hint txt-sm txt-right")},m(o,i){u(o,n,i)},d(o){o&&m(n)}}}function We(d){let n,o=d[12]+"",i,f=d[14]<d[4].length-1?", ":"",h;return{c(){n=e("code"),i=_(o),h=_(f)},m(r,b){u(r,n,b),t(n,i),u(r,h,b)},p(r,b){b&16&&o!==(o=r[12]+"")&&ke(i,o),b&16&&f!==(f=r[14]<r[4].length-1?", ":"")&&ke(h,f)},d(r){r&&(m(n),m(h))}}}function Xe(d,n){let o,i,f;function h(){return n[6](n[7])}return{key:d,first:null,c(){o=e("button"),o.textContent=`${n[7].code} `,a(o,"type","button"),a(o,"class","tab-item"),Le(o,"active",n[2]===n[7].code),this.first=o},m(r,b){u(r,o,b),i||(f=ll(o,"click",h),i=!0)},p(r,b){n=r,b&36&&Le(o,"active",n[2]===n[7].code)},d(r){r&&m(o),i=!1,f()}}}function Ye(d,n){let o,i,f,h;return i=new Fe({props:{content:n[7].body}}),{key:d,first:null,c(){o=e("div"),Qt(i.$$.fragment),f=s(),a(o,"class","tab-item"),Le(o,"active",n[2]===n[7].code),this.first=o},m(r,b){u(r,o,b),Ut(i,o,null),t(o,f),h=!0},p(r,b){n=r,(!h||b&36)&&Le(o,"active",n[2]===n[7].code)},i(r){h||($t(i.$$.fragment,r),h=!0)},o(r){Ct(i.$$.fragment,r),h=!1},d(r){r&&m(o),jt(i)}}}function gl(d){var Pe,Re,Oe,Se,Ee,Ne;let n,o,i=d[0].name+"",f,h,r,b,$,C,g,p=d[0].name+"",tt,kt,zt,E,Kt,H,rt,R,et,ne,Q,U,oe,ct,yt=d[0].name+"",lt,vt,ae,dt,pt,st,N,Jt,Ft,y,nt,Lt,Vt,At,j,ot,Tt,Wt,Pt,F,ft,Rt,ie,ut,re,M,Ot,at,St,O,mt,ce,z,Et,Xt,Nt,de,q,Yt,K,ht,pe,I,fe,B,ue,P,qt,J,bt,me,gt,he,x,Dt,it,Ht,be,Mt,Zt,V,_t,ge,It,_e,wt,we,W,G,xe,xt,te,X,ee,L,Y,S=[],Bt=new Map,$e,Z,v=[],Gt=new Map,A;E=new dl({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${d[3]}');

        ...

        // fetch a paginated records list
        const resultList = await pb.collection('${(Pe=d[0])==null?void 0:Pe.name}').getList(1, 50, {
            filter: 'created >= "2022-01-01 00:00:00" && someField1 != someField2',
        });

        // you can also fetch all records at once via getFullList
        const records = await pb.collection('${(Re=d[0])==null?void 0:Re.name}').getFullList({
            sort: '-created',
        });

        // or fetch only the first record that matches the specified filter
        const record = await pb.collection('${(Oe=d[0])==null?void 0:Oe.name}').getFirstListItem('someField="test"', {
            expand: 'relField1,relField2.subRelField',
        });
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${d[3]}');

        ...

        // fetch a paginated records list
        final resultList = await pb.collection('${(Se=d[0])==null?void 0:Se.name}').getList(
          page: 1,
          perPage: 50,
          filter: 'created >= "2022-01-01 00:00:00" && someField1 != someField2',
        );

        // you can also fetch all records at once via getFullList
        final records = await pb.collection('${(Ee=d[0])==null?void 0:Ee.name}').getFullList(
          sort: '-created',
        );

        // or fetch only the first record that matches the specified filter
        final record = await pb.collection('${(Ne=d[0])==null?void 0:Ne.name}').getFirstListItem(
          'someField="test"',
          expand: 'relField1,relField2.subRelField',
        );
    `}});let k=d[1]&&Ve();at=new Fe({props:{content:`
                        // DESC by created and ASC by id
                        ?sort=-created,id
                    `}});let le=se(d[4]),T=[];for(let l=0;l<le.length;l+=1)T[l]=We(Je(d,le,l));B=new Fe({props:{content:`
                        ?filter=(id='abc' && created>'2022-01-01')
                    `}}),P=new bl({}),it=new Fe({props:{content:"?expand=relField1,relField2.subRelField"}}),G=new pl({});let ye=se(d[5]);const Ae=l=>l[7].code;for(let l=0;l<ye.length;l+=1){let c=Ke(d,ye,l),w=Ae(c);Bt.set(w,S[l]=Xe(w,c))}let Ce=se(d[5]);const Te=l=>l[7].code;for(let l=0;l<Ce.length;l+=1){let c=ze(d,Ce,l),w=Te(c);Gt.set(w,v[l]=Ye(w,c))}return{c(){n=e("h3"),o=_("List/Search ("),f=_(i),h=_(")"),r=s(),b=e("div"),$=e("p"),C=_("Fetch a paginated "),g=e("strong"),tt=_(p),kt=_(" records list, supporting sorting and filtering."),zt=s(),Qt(E.$$.fragment),Kt=s(),H=e("h6"),H.textContent="API details",rt=s(),R=e("div"),et=e("strong"),et.textContent="GET",ne=s(),Q=e("div"),U=e("p"),oe=_("/api/collections/"),ct=e("strong"),lt=_(yt),vt=_("/records"),ae=s(),k&&k.c(),dt=s(),pt=e("div"),pt.textContent="Query parameters",st=s(),N=e("table"),Jt=e("thead"),Jt.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr>',Ft=s(),y=e("tbody"),nt=e("tr"),nt.innerHTML='<td>page</td> <td><span class="label">Number</span></td> <td>The page (aka. offset) of the paginated list (default to 1).</td>',Lt=s(),Vt=e("tr"),Vt.innerHTML='<td>perPage</td> <td><span class="label">Number</span></td> <td>Specify the max returned records per page (default to 30).</td>',At=s(),j=e("tr"),ot=e("td"),ot.textContent="sort",Tt=s(),Wt=e("td"),Wt.innerHTML='<span class="label">String</span>',Pt=s(),F=e("td"),ft=_("Specify the records order attribute(s). "),Rt=e("br"),ie=_(`
                Add `),ut=e("code"),ut.textContent="-",re=_(" / "),M=e("code"),M.textContent="+",Ot=_(` (default) in front of the attribute for DESC / ASC order.
                Ex.:
                `),Qt(at.$$.fragment),St=s(),O=e("p"),mt=e("strong"),mt.textContent="Supported record sort fields:",ce=s(),z=e("br"),Et=s(),Xt=e("code"),Xt.textContent="@random",Nt=_(`,
                    `);for(let l=0;l<T.length;l+=1)T[l].c();de=s(),q=e("tr"),Yt=e("td"),Yt.textContent="filter",K=s(),ht=e("td"),ht.innerHTML='<span class="label">String</span>',pe=s(),I=e("td"),fe=_(`Filter the returned records. Ex.:
                `),Qt(B.$$.fragment),ue=s(),Qt(P.$$.fragment),qt=s(),J=e("tr"),bt=e("td"),bt.textContent="expand",me=s(),gt=e("td"),gt.innerHTML='<span class="label">String</span>',he=s(),x=e("td"),Dt=_(`Auto expand record relations. Ex.:
                `),Qt(it.$$.fragment),Ht=_(`
                Supports up to 6-levels depth nested relations expansion. `),be=e("br"),Mt=_(`
                The expanded relations will be appended to each individual record under the
                `),Zt=e("code"),Zt.textContent="expand",V=_(" property (eg. "),_t=e("code"),_t.textContent='"expand": {"relField1": {...}, ...}',ge=_(`).
                `),It=e("br"),_e=_(`
                Only the relations to which the request user has permissions to `),wt=e("strong"),wt.textContent="view",we=_(" will be expanded."),W=s(),Qt(G.$$.fragment),xe=s(),xt=e("tr"),xt.innerHTML=`<td id="query-page">skipTotal</td> <td><span class="label">Boolean</span></td> <td>If it is set the total counts query will be skipped and the response fields
                <code>totalItems</code> and <code>totalPages</code> will have <code>-1</code> value.
                <br/>
                This could drastically speed up the search queries when the total counters are not needed or cursor
                based pagination is used.
                <br/>
                For optimization purposes, it is set by default for the
                <code>getFirstListItem()</code>
                and
                <code>getFullList()</code> SDKs methods.</td>`,te=s(),X=e("div"),X.textContent="Responses",ee=s(),L=e("div"),Y=e("div");for(let l=0;l<S.length;l+=1)S[l].c();$e=s(),Z=e("div");for(let l=0;l<v.length;l+=1)v[l].c();a(n,"class","m-b-sm"),a(b,"class","content txt-lg m-b-sm"),a(H,"class","m-b-xs"),a(et,"class","label label-primary"),a(Q,"class","content"),a(R,"class","alert alert-info"),a(pt,"class","section-title"),a(N,"class","table-compact table-border m-b-base"),a(X,"class","section-title"),a(Y,"class","tabs-header compact combined left"),a(Z,"class","tabs-content"),a(L,"class","tabs")},m(l,c){u(l,n,c),t(n,o),t(n,f),t(n,h),u(l,r,c),u(l,b,c),t(b,$),t($,C),t($,g),t(g,tt),t($,kt),u(l,zt,c),Ut(E,l,c),u(l,Kt,c),u(l,H,c),u(l,rt,c),u(l,R,c),t(R,et),t(R,ne),t(R,Q),t(Q,U),t(U,oe),t(U,ct),t(ct,lt),t(U,vt),t(R,ae),k&&k.m(R,null),u(l,dt,c),u(l,pt,c),u(l,st,c),u(l,N,c),t(N,Jt),t(N,Ft),t(N,y),t(y,nt),t(y,Lt),t(y,Vt),t(y,At),t(y,j),t(j,ot),t(j,Tt),t(j,Wt),t(j,Pt),t(j,F),t(F,ft),t(F,Rt),t(F,ie),t(F,ut),t(F,re),t(F,M),t(F,Ot),Ut(at,F,null),t(F,St),t(F,O),t(O,mt),t(O,ce),t(O,z),t(O,Et),t(O,Xt),t(O,Nt);for(let w=0;w<T.length;w+=1)T[w]&&T[w].m(O,null);t(y,de),t(y,q),t(q,Yt),t(q,K),t(q,ht),t(q,pe),t(q,I),t(I,fe),Ut(B,I,null),t(I,ue),Ut(P,I,null),t(y,qt),t(y,J),t(J,bt),t(J,me),t(J,gt),t(J,he),t(J,x),t(x,Dt),Ut(it,x,null),t(x,Ht),t(x,be),t(x,Mt),t(x,Zt),t(x,V),t(x,_t),t(x,ge),t(x,It),t(x,_e),t(x,wt),t(x,we),t(y,W),Ut(G,y,null),t(y,xe),t(y,xt),u(l,te,c),u(l,X,c),u(l,ee,c),u(l,L,c),t(L,Y);for(let w=0;w<S.length;w+=1)S[w]&&S[w].m(Y,null);t(L,$e),t(L,Z);for(let w=0;w<v.length;w+=1)v[w]&&v[w].m(Z,null);A=!0},p(l,[c]){var qe,De,He,Me,Ie,Be;(!A||c&1)&&i!==(i=l[0].name+"")&&ke(f,i),(!A||c&1)&&p!==(p=l[0].name+"")&&ke(tt,p);const w={};if(c&9&&(w.js=`
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
    `),E.$set(w),(!A||c&1)&&yt!==(yt=l[0].name+"")&&ke(lt,yt),l[1]?k||(k=Ve(),k.c(),k.m(R,null)):k&&(k.d(1),k=null),c&16){le=se(l[4]);let D;for(D=0;D<le.length;D+=1){const Ge=Je(l,le,D);T[D]?T[D].p(Ge,c):(T[D]=We(Ge),T[D].c(),T[D].m(O,null))}for(;D<T.length;D+=1)T[D].d(1);T.length=le.length}c&36&&(ye=se(l[5]),S=Ue(S,c,Ae,1,l,ye,Bt,Y,nl,Xe,null,Ke)),c&36&&(Ce=se(l[5]),ol(),v=Ue(v,c,Te,1,l,Ce,Gt,Z,al,Ye,null,ze),il())},i(l){if(!A){$t(E.$$.fragment,l),$t(at.$$.fragment,l),$t(B.$$.fragment,l),$t(P.$$.fragment,l),$t(it.$$.fragment,l),$t(G.$$.fragment,l);for(let c=0;c<Ce.length;c+=1)$t(v[c]);A=!0}},o(l){Ct(E.$$.fragment,l),Ct(at.$$.fragment,l),Ct(B.$$.fragment,l),Ct(P.$$.fragment,l),Ct(it.$$.fragment,l),Ct(G.$$.fragment,l);for(let c=0;c<v.length;c+=1)Ct(v[c]);A=!1},d(l){l&&(m(n),m(r),m(b),m(zt),m(Kt),m(H),m(rt),m(R),m(dt),m(pt),m(st),m(N),m(te),m(X),m(ee),m(L)),jt(E,l),k&&k.d(),jt(at),rl(T,l),jt(B),jt(P),jt(it),jt(G);for(let c=0;c<S.length;c+=1)S[c].d();for(let c=0;c<v.length;c+=1)v[c].d()}}}function _l(d,n,o){let i,f,h,{collection:r}=n,b=200,$=[];const C=g=>o(2,b=g.code);return d.$$set=g=>{"collection"in g&&o(0,r=g.collection)},d.$$.update=()=>{d.$$.dirty&1&&o(4,i=ve.getAllCollectionIdentifiers(r)),d.$$.dirty&1&&o(1,f=(r==null?void 0:r.listRule)===null),d.$$.dirty&3&&r!=null&&r.id&&($.push({code:200,body:JSON.stringify({page:1,perPage:30,totalPages:1,totalItems:2,items:[ve.dummyCollectionRecord(r),ve.dummyCollectionRecord(r)]},null,2)}),$.push({code:400,body:`
                {
                  "code": 400,
                  "message": "Something went wrong while processing your request. Invalid filter.",
                  "data": {}
                }
            `}),f&&$.push({code:403,body:`
                    {
                      "code": 403,
                      "message": "Only admins can access this action.",
                      "data": {}
                    }
                `}))},o(3,h=ve.getApiExampleUrl(cl.baseUrl)),[r,f,b,h,i,$,C]}class Cl extends Ze{constructor(n){super(),tl(this,n,_l,gl,el,{collection:0})}}export{Cl as default};
