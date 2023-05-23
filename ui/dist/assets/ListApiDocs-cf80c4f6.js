import{S as We,i as Xe,s as Ye,e,b as s,E as tl,f as i,g as u,u as Ze,y as Ie,o as m,w,h as t,M as ke,c as Zt,m as te,x as ge,N as Be,P as el,k as ll,Q as sl,n as nl,t as Bt,a as Gt,d as ee,R as ol,T as il,C as $e,p as al,r as ye}from"./index-a65ca895.js";import{S as rl}from"./SdkTabs-ad912c8f.js";function cl(c){let n,o,a;return{c(){n=e("span"),n.textContent="Show details",o=s(),a=e("i"),i(n,"class","txt"),i(a,"class","ri-arrow-down-s-line")},m(p,b){u(p,n,b),u(p,o,b),u(p,a,b)},d(p){p&&m(n),p&&m(o),p&&m(a)}}}function dl(c){let n,o,a;return{c(){n=e("span"),n.textContent="Hide details",o=s(),a=e("i"),i(n,"class","txt"),i(a,"class","ri-arrow-up-s-line")},m(p,b){u(p,n,b),u(p,o,b),u(p,a,b)},d(p){p&&m(n),p&&m(o),p&&m(a)}}}function Ge(c){let n,o,a,p,b,d,h,g,x,_,f,Z,Ct,Ut,O,jt,H,at,R,tt,le,U,j,se,rt,$t,et,kt,ne,ct,dt,lt,E,Jt,yt,L,st,vt,Qt,Ft,J,nt,Lt,zt,At,T,ft,Tt,oe,pt,ie,D,Pt,ot,Rt,S,ut,ae,Q,St,Kt,Ot,re,N,Vt,z,mt,ce,I,de,B,fe,P,Et,K,bt,pe,ht,ue,$,Nt,it,qt,me,Mt,Wt,V,_t,be,Ht,he,xt,_e,W,wt,Xt,X,Yt,q,gt,y,Dt,xe,Y,A,It,G,v;return{c(){n=e("p"),n.innerHTML=`The syntax basically follows the format
        <code><span class="txt-success">OPERAND</span> 
            <span class="txt-danger">OPERATOR</span> 
            <span class="txt-success">OPERAND</span></code>, where:`,o=s(),a=e("ul"),p=e("li"),p.innerHTML=`<code class="txt-success">OPERAND</code> - could be any of the above field literal, string (single
            or double quoted), number, null, true, false`,b=s(),d=e("li"),h=e("code"),h.textContent="OPERATOR",g=w(` - is one of:
            `),x=e("br"),_=s(),f=e("ul"),Z=e("li"),Ct=e("code"),Ct.textContent="=",Ut=s(),O=e("span"),O.textContent="Equal",jt=s(),H=e("li"),at=e("code"),at.textContent="!=",R=s(),tt=e("span"),tt.textContent="NOT equal",le=s(),U=e("li"),j=e("code"),j.textContent=">",se=s(),rt=e("span"),rt.textContent="Greater than",$t=s(),et=e("li"),kt=e("code"),kt.textContent=">=",ne=s(),ct=e("span"),ct.textContent="Greater than or equal",dt=s(),lt=e("li"),E=e("code"),E.textContent="<",Jt=s(),yt=e("span"),yt.textContent="Less than",L=s(),st=e("li"),vt=e("code"),vt.textContent="<=",Qt=s(),Ft=e("span"),Ft.textContent="Less than or equal",J=s(),nt=e("li"),Lt=e("code"),Lt.textContent="~",zt=s(),At=e("span"),At.textContent=`Like/Contains (if not specified auto wraps the right string OPERAND in a "%" for
                        wildcard match)`,T=s(),ft=e("li"),Tt=e("code"),Tt.textContent="!~",oe=s(),pt=e("span"),pt.textContent=`NOT Like/Contains (if not specified auto wraps the right string OPERAND in a "%" for
                        wildcard match)`,ie=s(),D=e("li"),Pt=e("code"),Pt.textContent="?=",ot=s(),Rt=e("em"),Rt.textContent="Any/At least one of",S=s(),ut=e("span"),ut.textContent="Equal",ae=s(),Q=e("li"),St=e("code"),St.textContent="?!=",Kt=s(),Ot=e("em"),Ot.textContent="Any/At least one of",re=s(),N=e("span"),N.textContent="NOT equal",Vt=s(),z=e("li"),mt=e("code"),mt.textContent="?>",ce=s(),I=e("em"),I.textContent="Any/At least one of",de=s(),B=e("span"),B.textContent="Greater than",fe=s(),P=e("li"),Et=e("code"),Et.textContent="?>=",K=s(),bt=e("em"),bt.textContent="Any/At least one of",pe=s(),ht=e("span"),ht.textContent="Greater than or equal",ue=s(),$=e("li"),Nt=e("code"),Nt.textContent="?<",it=s(),qt=e("em"),qt.textContent="Any/At least one of",me=s(),Mt=e("span"),Mt.textContent="Less than",Wt=s(),V=e("li"),_t=e("code"),_t.textContent="?<=",be=s(),Ht=e("em"),Ht.textContent="Any/At least one of",he=s(),xt=e("span"),xt.textContent="Less than or equal",_e=s(),W=e("li"),wt=e("code"),wt.textContent="?~",Xt=s(),X=e("em"),X.textContent="Any/At least one of",Yt=s(),q=e("span"),q.textContent=`Like/Contains (if not specified auto wraps the right string OPERAND in a "%" for
                        wildcard match)`,gt=s(),y=e("li"),Dt=e("code"),Dt.textContent="?!~",xe=s(),Y=e("em"),Y.textContent="Any/At least one of",A=s(),It=e("span"),It.textContent=`NOT Like/Contains (if not specified auto wraps the right string OPERAND in a "%" for
                        wildcard match)`,G=s(),v=e("p"),v.innerHTML=`To group and combine several expressions you could use brackets
        <code>(...)</code>, <code>&amp;&amp;</code> (AND) and <code>||</code> (OR) tokens.`,i(h,"class","txt-danger"),i(Ct,"class","filter-op svelte-1w7s5nw"),i(O,"class","txt"),i(at,"class","filter-op svelte-1w7s5nw"),i(tt,"class","txt"),i(j,"class","filter-op svelte-1w7s5nw"),i(rt,"class","txt"),i(kt,"class","filter-op svelte-1w7s5nw"),i(ct,"class","txt"),i(E,"class","filter-op svelte-1w7s5nw"),i(yt,"class","txt"),i(vt,"class","filter-op svelte-1w7s5nw"),i(Ft,"class","txt"),i(Lt,"class","filter-op svelte-1w7s5nw"),i(At,"class","txt"),i(Tt,"class","filter-op svelte-1w7s5nw"),i(pt,"class","txt"),i(Pt,"class","filter-op svelte-1w7s5nw"),i(Rt,"class","txt-hint"),i(ut,"class","txt"),i(St,"class","filter-op svelte-1w7s5nw"),i(Ot,"class","txt-hint"),i(N,"class","txt"),i(mt,"class","filter-op svelte-1w7s5nw"),i(I,"class","txt-hint"),i(B,"class","txt"),i(Et,"class","filter-op svelte-1w7s5nw"),i(bt,"class","txt-hint"),i(ht,"class","txt"),i(Nt,"class","filter-op svelte-1w7s5nw"),i(qt,"class","txt-hint"),i(Mt,"class","txt"),i(_t,"class","filter-op svelte-1w7s5nw"),i(Ht,"class","txt-hint"),i(xt,"class","txt"),i(wt,"class","filter-op svelte-1w7s5nw"),i(X,"class","txt-hint"),i(q,"class","txt"),i(Dt,"class","filter-op svelte-1w7s5nw"),i(Y,"class","txt-hint"),i(It,"class","txt")},m(F,k){u(F,n,k),u(F,o,k),u(F,a,k),t(a,p),t(a,b),t(a,d),t(d,h),t(d,g),t(d,x),t(d,_),t(d,f),t(f,Z),t(Z,Ct),t(Z,Ut),t(Z,O),t(f,jt),t(f,H),t(H,at),t(H,R),t(H,tt),t(f,le),t(f,U),t(U,j),t(U,se),t(U,rt),t(f,$t),t(f,et),t(et,kt),t(et,ne),t(et,ct),t(f,dt),t(f,lt),t(lt,E),t(lt,Jt),t(lt,yt),t(f,L),t(f,st),t(st,vt),t(st,Qt),t(st,Ft),t(f,J),t(f,nt),t(nt,Lt),t(nt,zt),t(nt,At),t(f,T),t(f,ft),t(ft,Tt),t(ft,oe),t(ft,pt),t(f,ie),t(f,D),t(D,Pt),t(D,ot),t(D,Rt),t(D,S),t(D,ut),t(f,ae),t(f,Q),t(Q,St),t(Q,Kt),t(Q,Ot),t(Q,re),t(Q,N),t(f,Vt),t(f,z),t(z,mt),t(z,ce),t(z,I),t(z,de),t(z,B),t(f,fe),t(f,P),t(P,Et),t(P,K),t(P,bt),t(P,pe),t(P,ht),t(f,ue),t(f,$),t($,Nt),t($,it),t($,qt),t($,me),t($,Mt),t(f,Wt),t(f,V),t(V,_t),t(V,be),t(V,Ht),t(V,he),t(V,xt),t(f,_e),t(f,W),t(W,wt),t(W,Xt),t(W,X),t(W,Yt),t(W,q),t(f,gt),t(f,y),t(y,Dt),t(y,xe),t(y,Y),t(y,A),t(y,It),u(F,G,k),u(F,v,k)},d(F){F&&m(n),F&&m(o),F&&m(a),F&&m(G),F&&m(v)}}}function fl(c){let n,o,a,p,b;function d(_,f){return _[0]?dl:cl}let h=d(c),g=h(c),x=c[0]&&Ge();return{c(){n=e("button"),g.c(),o=s(),x&&x.c(),a=tl(),i(n,"class","btn btn-sm btn-secondary m-t-10")},m(_,f){u(_,n,f),g.m(n,null),u(_,o,f),x&&x.m(_,f),u(_,a,f),p||(b=Ze(n,"click",c[1]),p=!0)},p(_,[f]){h!==(h=d(_))&&(g.d(1),g=h(_),g&&(g.c(),g.m(n,null))),_[0]?x||(x=Ge(),x.c(),x.m(a.parentNode,a)):x&&(x.d(1),x=null)},i:Ie,o:Ie,d(_){_&&m(n),g.d(),_&&m(o),x&&x.d(_),_&&m(a),p=!1,b()}}}function pl(c,n,o){let a=!1;function p(){o(0,a=!a)}return[a,p]}class ul extends We{constructor(n){super(),Xe(this,n,pl,fl,Ye,{})}}function Ue(c,n,o){const a=c.slice();return a[7]=n[o],a}function je(c,n,o){const a=c.slice();return a[7]=n[o],a}function Je(c,n,o){const a=c.slice();return a[12]=n[o],a[14]=o,a}function Qe(c){let n;return{c(){n=e("p"),n.innerHTML="Requires admin <code>Authorization:TOKEN</code> header",i(n,"class","txt-hint txt-sm txt-right")},m(o,a){u(o,n,a)},d(o){o&&m(n)}}}function ze(c){let n,o=c[12]+"",a,p=c[14]<c[4].length-1?", ":"",b;return{c(){n=e("code"),a=w(o),b=w(p)},m(d,h){u(d,n,h),t(n,a),u(d,b,h)},p(d,h){h&16&&o!==(o=d[12]+"")&&ge(a,o),h&16&&p!==(p=d[14]<d[4].length-1?", ":"")&&ge(b,p)},d(d){d&&m(n),d&&m(b)}}}function Ke(c,n){let o,a=n[7].code+"",p,b,d,h;function g(){return n[6](n[7])}return{key:c,first:null,c(){o=e("button"),p=w(a),b=s(),i(o,"type","button"),i(o,"class","tab-item"),ye(o,"active",n[2]===n[7].code),this.first=o},m(x,_){u(x,o,_),t(o,p),t(o,b),d||(h=Ze(o,"click",g),d=!0)},p(x,_){n=x,_&36&&ye(o,"active",n[2]===n[7].code)},d(x){x&&m(o),d=!1,h()}}}function Ve(c,n){let o,a,p,b;return a=new ke({props:{content:n[7].body}}),{key:c,first:null,c(){o=e("div"),Zt(a.$$.fragment),p=s(),i(o,"class","tab-item"),ye(o,"active",n[2]===n[7].code),this.first=o},m(d,h){u(d,o,h),te(a,o,null),t(o,p),b=!0},p(d,h){n=d,(!b||h&36)&&ye(o,"active",n[2]===n[7].code)},i(d){b||(Bt(a.$$.fragment,d),b=!0)},o(d){Gt(a.$$.fragment,d),b=!1},d(d){d&&m(o),ee(a)}}}function ml(c){var Le,Ae,Te,Pe,Re,Se;let n,o,a=c[0].name+"",p,b,d,h,g,x,_,f=c[0].name+"",Z,Ct,Ut,O,jt,H,at,R,tt,le,U,j,se,rt,$t=c[0].name+"",et,kt,ne,ct,dt,lt,E,Jt,yt,L,st,vt,Qt,Ft,J,nt,Lt,zt,At,T,ft,Tt,oe,pt,ie,D,Pt,ot,Rt,S,ut,ae,Q,St,Kt,Ot,re,N,Vt,z,mt,ce,I,de,B,fe,P,Et,K,bt,pe,ht,ue,$,Nt,it,qt,me,Mt,Wt,V,_t,be,Ht,he,xt,_e,W,wt,Xt,X,Yt,q,gt,y=[],Dt=new Map,xe,Y,A=[],It=new Map,G;O=new rl({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${c[3]}');

        ...

        // fetch a paginated records list
        const resultList = await pb.collection('${(Le=c[0])==null?void 0:Le.name}').getList(1, 50, {
            filter: 'created >= "2022-01-01 00:00:00" && someField1 != someField2',
        });

        // you can also fetch all records at once via getFullList
        const records = await pb.collection('${(Ae=c[0])==null?void 0:Ae.name}').getFullList({
            sort: '-created',
        });

        // or fetch only the first record that matches the specified filter
        const record = await pb.collection('${(Te=c[0])==null?void 0:Te.name}').getFirstListItem('someField="test"', {
            expand: 'relField1,relField2.subRelField',
        });
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${c[3]}');

        ...

        // fetch a paginated records list
        final resultList = await pb.collection('${(Pe=c[0])==null?void 0:Pe.name}').getList(
          page: 1,
          perPage: 50,
          filter: 'created >= "2022-01-01 00:00:00" && someField1 != someField2',
        );

        // you can also fetch all records at once via getFullList
        final records = await pb.collection('${(Re=c[0])==null?void 0:Re.name}').getFullList(
          sort: '-created',
        );

        // or fetch only the first record that matches the specified filter
        final record = await pb.collection('${(Se=c[0])==null?void 0:Se.name}').getFirstListItem(
          'someField="test"',
          expand: 'relField1,relField2.subRelField',
        );
    `}});let v=c[1]&&Qe();ot=new ke({props:{content:`
                        // DESC by created and ASC by id
                        ?sort=-created,id
                    `}});let F=c[4],k=[];for(let l=0;l<F.length;l+=1)k[l]=ze(Je(c,F,l));B=new ke({props:{content:`
                        ?filter=(id='abc' && created>'2022-01-01')
                    `}}),P=new ul({}),it=new ke({props:{content:"?expand=relField1,relField2.subRelField"}});let Ce=c[5];const ve=l=>l[7].code;for(let l=0;l<Ce.length;l+=1){let r=je(c,Ce,l),C=ve(r);Dt.set(C,y[l]=Ke(C,r))}let we=c[5];const Fe=l=>l[7].code;for(let l=0;l<we.length;l+=1){let r=Ue(c,we,l),C=Fe(r);It.set(C,A[l]=Ve(C,r))}return{c(){n=e("h3"),o=w("List/Search ("),p=w(a),b=w(")"),d=s(),h=e("div"),g=e("p"),x=w("Fetch a paginated "),_=e("strong"),Z=w(f),Ct=w(" records list, supporting sorting and filtering."),Ut=s(),Zt(O.$$.fragment),jt=s(),H=e("h6"),H.textContent="API details",at=s(),R=e("div"),tt=e("strong"),tt.textContent="GET",le=s(),U=e("div"),j=e("p"),se=w("/api/collections/"),rt=e("strong"),et=w($t),kt=w("/records"),ne=s(),v&&v.c(),ct=s(),dt=e("div"),dt.textContent="Query parameters",lt=s(),E=e("table"),Jt=e("thead"),Jt.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr>`,yt=s(),L=e("tbody"),st=e("tr"),st.innerHTML=`<td>page</td> 
            <td><span class="label">Number</span></td> 
            <td>The page (aka. offset) of the paginated list (default to 1).</td>`,vt=s(),Qt=e("tr"),Qt.innerHTML=`<td>perPage</td> 
            <td><span class="label">Number</span></td> 
            <td>Specify the max returned records per page (default to 30).</td>`,Ft=s(),J=e("tr"),nt=e("td"),nt.textContent="sort",Lt=s(),zt=e("td"),zt.innerHTML='<span class="label">String</span>',At=s(),T=e("td"),ft=w("Specify the records order attribute(s). "),Tt=e("br"),oe=w(`
                Add `),pt=e("code"),pt.textContent="-",ie=w(" / "),D=e("code"),D.textContent="+",Pt=w(` (default) in front of the attribute for DESC / ASC order.
                Ex.:
                `),Zt(ot.$$.fragment),Rt=s(),S=e("p"),ut=e("strong"),ut.textContent="Supported record sort fields:",ae=s(),Q=e("br"),St=s(),Kt=e("code"),Kt.textContent="@random",Ot=w(`,
                    `);for(let l=0;l<k.length;l+=1)k[l].c();re=s(),N=e("tr"),Vt=e("td"),Vt.textContent="filter",z=s(),mt=e("td"),mt.innerHTML='<span class="label">String</span>',ce=s(),I=e("td"),de=w(`Filter the returned records. Ex.:
                `),Zt(B.$$.fragment),fe=s(),Zt(P.$$.fragment),Et=s(),K=e("tr"),bt=e("td"),bt.textContent="expand",pe=s(),ht=e("td"),ht.innerHTML='<span class="label">String</span>',ue=s(),$=e("td"),Nt=w(`Auto expand record relations. Ex.:
                `),Zt(it.$$.fragment),qt=w(`
                Supports up to 6-levels depth nested relations expansion. `),me=e("br"),Mt=w(`
                The expanded relations will be appended to each individual record under the
                `),Wt=e("code"),Wt.textContent="expand",V=w(" property (eg. "),_t=e("code"),_t.textContent='"expand": {"relField1": {...}, ...}',be=w(`).
                `),Ht=e("br"),he=w(`
                Only the relations to which the request user has permissions to `),xt=e("strong"),xt.textContent="view",_e=w(" will be expanded."),W=s(),wt=e("tr"),wt.innerHTML=`<td id="query-page">fields</td> 
            <td><span class="label">String</span></td> 
            <td>Comma separated string of the fields to return in the JSON response
                <em>(by default returns all fields)</em>.</td>`,Xt=s(),X=e("div"),X.textContent="Responses",Yt=s(),q=e("div"),gt=e("div");for(let l=0;l<y.length;l+=1)y[l].c();xe=s(),Y=e("div");for(let l=0;l<A.length;l+=1)A[l].c();i(n,"class","m-b-sm"),i(h,"class","content txt-lg m-b-sm"),i(H,"class","m-b-xs"),i(tt,"class","label label-primary"),i(U,"class","content"),i(R,"class","alert alert-info"),i(dt,"class","section-title"),i(E,"class","table-compact table-border m-b-base"),i(X,"class","section-title"),i(gt,"class","tabs-header compact left"),i(Y,"class","tabs-content"),i(q,"class","tabs")},m(l,r){u(l,n,r),t(n,o),t(n,p),t(n,b),u(l,d,r),u(l,h,r),t(h,g),t(g,x),t(g,_),t(_,Z),t(g,Ct),u(l,Ut,r),te(O,l,r),u(l,jt,r),u(l,H,r),u(l,at,r),u(l,R,r),t(R,tt),t(R,le),t(R,U),t(U,j),t(j,se),t(j,rt),t(rt,et),t(j,kt),t(R,ne),v&&v.m(R,null),u(l,ct,r),u(l,dt,r),u(l,lt,r),u(l,E,r),t(E,Jt),t(E,yt),t(E,L),t(L,st),t(L,vt),t(L,Qt),t(L,Ft),t(L,J),t(J,nt),t(J,Lt),t(J,zt),t(J,At),t(J,T),t(T,ft),t(T,Tt),t(T,oe),t(T,pt),t(T,ie),t(T,D),t(T,Pt),te(ot,T,null),t(T,Rt),t(T,S),t(S,ut),t(S,ae),t(S,Q),t(S,St),t(S,Kt),t(S,Ot);for(let C=0;C<k.length;C+=1)k[C]&&k[C].m(S,null);t(L,re),t(L,N),t(N,Vt),t(N,z),t(N,mt),t(N,ce),t(N,I),t(I,de),te(B,I,null),t(I,fe),te(P,I,null),t(L,Et),t(L,K),t(K,bt),t(K,pe),t(K,ht),t(K,ue),t(K,$),t($,Nt),te(it,$,null),t($,qt),t($,me),t($,Mt),t($,Wt),t($,V),t($,_t),t($,be),t($,Ht),t($,he),t($,xt),t($,_e),t(L,W),t(L,wt),u(l,Xt,r),u(l,X,r),u(l,Yt,r),u(l,q,r),t(q,gt);for(let C=0;C<y.length;C+=1)y[C]&&y[C].m(gt,null);t(q,xe),t(q,Y);for(let C=0;C<A.length;C+=1)A[C]&&A[C].m(Y,null);G=!0},p(l,[r]){var Oe,Ee,Ne,qe,Me,He;(!G||r&1)&&a!==(a=l[0].name+"")&&ge(p,a),(!G||r&1)&&f!==(f=l[0].name+"")&&ge(Z,f);const C={};if(r&9&&(C.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${l[3]}');

        ...

        // fetch a paginated records list
        const resultList = await pb.collection('${(Oe=l[0])==null?void 0:Oe.name}').getList(1, 50, {
            filter: 'created >= "2022-01-01 00:00:00" && someField1 != someField2',
        });

        // you can also fetch all records at once via getFullList
        const records = await pb.collection('${(Ee=l[0])==null?void 0:Ee.name}').getFullList({
            sort: '-created',
        });

        // or fetch only the first record that matches the specified filter
        const record = await pb.collection('${(Ne=l[0])==null?void 0:Ne.name}').getFirstListItem('someField="test"', {
            expand: 'relField1,relField2.subRelField',
        });
    `),r&9&&(C.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${l[3]}');

        ...

        // fetch a paginated records list
        final resultList = await pb.collection('${(qe=l[0])==null?void 0:qe.name}').getList(
          page: 1,
          perPage: 50,
          filter: 'created >= "2022-01-01 00:00:00" && someField1 != someField2',
        );

        // you can also fetch all records at once via getFullList
        final records = await pb.collection('${(Me=l[0])==null?void 0:Me.name}').getFullList(
          sort: '-created',
        );

        // or fetch only the first record that matches the specified filter
        final record = await pb.collection('${(He=l[0])==null?void 0:He.name}').getFirstListItem(
          'someField="test"',
          expand: 'relField1,relField2.subRelField',
        );
    `),O.$set(C),(!G||r&1)&&$t!==($t=l[0].name+"")&&ge(et,$t),l[1]?v||(v=Qe(),v.c(),v.m(R,null)):v&&(v.d(1),v=null),r&16){F=l[4];let M;for(M=0;M<F.length;M+=1){const De=Je(l,F,M);k[M]?k[M].p(De,r):(k[M]=ze(De),k[M].c(),k[M].m(S,null))}for(;M<k.length;M+=1)k[M].d(1);k.length=F.length}r&36&&(Ce=l[5],y=Be(y,r,ve,1,l,Ce,Dt,gt,el,Ke,null,je)),r&36&&(we=l[5],ll(),A=Be(A,r,Fe,1,l,we,It,Y,sl,Ve,null,Ue),nl())},i(l){if(!G){Bt(O.$$.fragment,l),Bt(ot.$$.fragment,l),Bt(B.$$.fragment,l),Bt(P.$$.fragment,l),Bt(it.$$.fragment,l);for(let r=0;r<we.length;r+=1)Bt(A[r]);G=!0}},o(l){Gt(O.$$.fragment,l),Gt(ot.$$.fragment,l),Gt(B.$$.fragment,l),Gt(P.$$.fragment,l),Gt(it.$$.fragment,l);for(let r=0;r<A.length;r+=1)Gt(A[r]);G=!1},d(l){l&&m(n),l&&m(d),l&&m(h),l&&m(Ut),ee(O,l),l&&m(jt),l&&m(H),l&&m(at),l&&m(R),v&&v.d(),l&&m(ct),l&&m(dt),l&&m(lt),l&&m(E),ee(ot),ol(k,l),ee(B),ee(P),ee(it),l&&m(Xt),l&&m(X),l&&m(Yt),l&&m(q);for(let r=0;r<y.length;r+=1)y[r].d();for(let r=0;r<A.length;r+=1)A[r].d()}}}function bl(c,n,o){let a,p,b,{collection:d=new il}=n,h=200,g=[];const x=_=>o(2,h=_.code);return c.$$set=_=>{"collection"in _&&o(0,d=_.collection)},c.$$.update=()=>{c.$$.dirty&1&&o(4,a=$e.getAllCollectionIdentifiers(d)),c.$$.dirty&1&&o(1,p=(d==null?void 0:d.listRule)===null),c.$$.dirty&3&&d!=null&&d.id&&(g.push({code:200,body:JSON.stringify({page:1,perPage:30,totalPages:1,totalItems:2,items:[$e.dummyCollectionRecord(d),$e.dummyCollectionRecord(d)]},null,2)}),g.push({code:400,body:`
                {
                  "code": 400,
                  "message": "Something went wrong while processing your request. Invalid filter.",
                  "data": {}
                }
            `}),p&&g.push({code:403,body:`
                    {
                      "code": 403,
                      "message": "Only admins can access this action.",
                      "data": {}
                    }
                `}))},o(3,b=$e.getApiExampleUrl(al.baseUrl)),[d,p,h,b,a,g,x]}class xl extends We{constructor(n){super(),Xe(this,n,bl,ml,Ye,{collection:0})}}export{xl as default};
