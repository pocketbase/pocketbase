import{S as Ke,i as Ve,s as We,e,b as s,E as Ye,f as i,g as u,u as Xe,y as He,o as m,w,h as t,N as $e,c as Zt,m as te,x as Ce,O as Me,P as Ze,k as tl,Q as el,n as ll,t as Gt,a as Ut,d as ee,R as sl,T as nl,C as ge,p as ol,r as ke}from"./index-4eea3e34.js";import{S as il}from"./SdkTabs-5d6cc1d4.js";function al(c){let n,o,a;return{c(){n=e("span"),n.textContent="Show details",o=s(),a=e("i"),i(n,"class","txt"),i(a,"class","ri-arrow-down-s-line")},m(p,b){u(p,n,b),u(p,o,b),u(p,a,b)},d(p){p&&m(n),p&&m(o),p&&m(a)}}}function rl(c){let n,o,a;return{c(){n=e("span"),n.textContent="Hide details",o=s(),a=e("i"),i(n,"class","txt"),i(a,"class","ri-arrow-up-s-line")},m(p,b){u(p,n,b),u(p,o,b),u(p,a,b)},d(p){p&&m(n),p&&m(o),p&&m(a)}}}function Ie(c){let n,o,a,p,b,d,h,C,x,_,f,et,kt,jt,S,Qt,H,ct,R,lt,le,U,j,se,dt,vt,st,yt,ne,ft,pt,nt,E,zt,Ft,T,ot,Lt,Jt,At,Q,it,Tt,Kt,Pt,L,ut,Rt,oe,mt,ie,M,Ot,at,St,O,bt,ae,z,Et,Vt,Nt,re,N,Wt,J,ht,ce,I,de,B,fe,P,qt,K,_t,pe,xt,ue,$,Dt,rt,Ht,me,Mt,Xt,V,wt,be,It,he,Ct,_e,G,W,Yt,q,gt,A,xe,X,Y,y,Bt,Z,v,tt,k;return{c(){n=e("p"),n.innerHTML=`The syntax basically follows the format
        <code><span class="txt-success">OPERAND</span> 
            <span class="txt-danger">OPERATOR</span> 
            <span class="txt-success">OPERAND</span></code>, where:`,o=s(),a=e("ul"),p=e("li"),p.innerHTML=`<code class="txt-success">OPERAND</code> - could be any of the above field literal, string (single
            or double quoted), number, null, true, false`,b=s(),d=e("li"),h=e("code"),h.textContent="OPERATOR",C=w(` - is one of:
            `),x=e("br"),_=s(),f=e("ul"),et=e("li"),kt=e("code"),kt.textContent="=",jt=s(),S=e("span"),S.textContent="Equal",Qt=s(),H=e("li"),ct=e("code"),ct.textContent="!=",R=s(),lt=e("span"),lt.textContent="NOT equal",le=s(),U=e("li"),j=e("code"),j.textContent=">",se=s(),dt=e("span"),dt.textContent="Greater than",vt=s(),st=e("li"),yt=e("code"),yt.textContent=">=",ne=s(),ft=e("span"),ft.textContent="Greater than or equal",pt=s(),nt=e("li"),E=e("code"),E.textContent="<",zt=s(),Ft=e("span"),Ft.textContent="Less than",T=s(),ot=e("li"),Lt=e("code"),Lt.textContent="<=",Jt=s(),At=e("span"),At.textContent="Less than or equal",Q=s(),it=e("li"),Tt=e("code"),Tt.textContent="~",Kt=s(),Pt=e("span"),Pt.textContent=`Like/Contains (if not specified auto wraps the right string OPERAND in a "%" for
                        wildcard match)`,L=s(),ut=e("li"),Rt=e("code"),Rt.textContent="!~",oe=s(),mt=e("span"),mt.textContent=`NOT Like/Contains (if not specified auto wraps the right string OPERAND in a "%" for
                        wildcard match)`,ie=s(),M=e("li"),Ot=e("code"),Ot.textContent="?=",at=s(),St=e("em"),St.textContent="Any/At least one of",O=s(),bt=e("span"),bt.textContent="Equal",ae=s(),z=e("li"),Et=e("code"),Et.textContent="?!=",Vt=s(),Nt=e("em"),Nt.textContent="Any/At least one of",re=s(),N=e("span"),N.textContent="NOT equal",Wt=s(),J=e("li"),ht=e("code"),ht.textContent="?>",ce=s(),I=e("em"),I.textContent="Any/At least one of",de=s(),B=e("span"),B.textContent="Greater than",fe=s(),P=e("li"),qt=e("code"),qt.textContent="?>=",K=s(),_t=e("em"),_t.textContent="Any/At least one of",pe=s(),xt=e("span"),xt.textContent="Greater than or equal",ue=s(),$=e("li"),Dt=e("code"),Dt.textContent="?<",rt=s(),Ht=e("em"),Ht.textContent="Any/At least one of",me=s(),Mt=e("span"),Mt.textContent="Less than",Xt=s(),V=e("li"),wt=e("code"),wt.textContent="?<=",be=s(),It=e("em"),It.textContent="Any/At least one of",he=s(),Ct=e("span"),Ct.textContent="Less than or equal",_e=s(),G=e("li"),W=e("code"),W.textContent="?~",Yt=s(),q=e("em"),q.textContent="Any/At least one of",gt=s(),A=e("span"),A.textContent=`Like/Contains (if not specified auto wraps the right string OPERAND in a "%" for
                        wildcard match)`,xe=s(),X=e("li"),Y=e("code"),Y.textContent="?!~",y=s(),Bt=e("em"),Bt.textContent="Any/At least one of",Z=s(),v=e("span"),v.textContent=`NOT Like/Contains (if not specified auto wraps the right string OPERAND in a "%" for
                        wildcard match)`,tt=s(),k=e("p"),k.innerHTML=`To group and combine several expressions you could use brackets
        <code>(...)</code>, <code>&amp;&amp;</code> (AND) and <code>||</code> (OR) tokens.`,i(h,"class","txt-danger"),i(kt,"class","filter-op svelte-1w7s5nw"),i(S,"class","txt"),i(ct,"class","filter-op svelte-1w7s5nw"),i(lt,"class","txt"),i(j,"class","filter-op svelte-1w7s5nw"),i(dt,"class","txt"),i(yt,"class","filter-op svelte-1w7s5nw"),i(ft,"class","txt"),i(E,"class","filter-op svelte-1w7s5nw"),i(Ft,"class","txt"),i(Lt,"class","filter-op svelte-1w7s5nw"),i(At,"class","txt"),i(Tt,"class","filter-op svelte-1w7s5nw"),i(Pt,"class","txt"),i(Rt,"class","filter-op svelte-1w7s5nw"),i(mt,"class","txt"),i(Ot,"class","filter-op svelte-1w7s5nw"),i(St,"class","txt-hint"),i(bt,"class","txt"),i(Et,"class","filter-op svelte-1w7s5nw"),i(Nt,"class","txt-hint"),i(N,"class","txt"),i(ht,"class","filter-op svelte-1w7s5nw"),i(I,"class","txt-hint"),i(B,"class","txt"),i(qt,"class","filter-op svelte-1w7s5nw"),i(_t,"class","txt-hint"),i(xt,"class","txt"),i(Dt,"class","filter-op svelte-1w7s5nw"),i(Ht,"class","txt-hint"),i(Mt,"class","txt"),i(wt,"class","filter-op svelte-1w7s5nw"),i(It,"class","txt-hint"),i(Ct,"class","txt"),i(W,"class","filter-op svelte-1w7s5nw"),i(q,"class","txt-hint"),i(A,"class","txt"),i(Y,"class","filter-op svelte-1w7s5nw"),i(Bt,"class","txt-hint"),i(v,"class","txt")},m(F,$t){u(F,n,$t),u(F,o,$t),u(F,a,$t),t(a,p),t(a,b),t(a,d),t(d,h),t(d,C),t(d,x),t(d,_),t(d,f),t(f,et),t(et,kt),t(et,jt),t(et,S),t(f,Qt),t(f,H),t(H,ct),t(H,R),t(H,lt),t(f,le),t(f,U),t(U,j),t(U,se),t(U,dt),t(f,vt),t(f,st),t(st,yt),t(st,ne),t(st,ft),t(f,pt),t(f,nt),t(nt,E),t(nt,zt),t(nt,Ft),t(f,T),t(f,ot),t(ot,Lt),t(ot,Jt),t(ot,At),t(f,Q),t(f,it),t(it,Tt),t(it,Kt),t(it,Pt),t(f,L),t(f,ut),t(ut,Rt),t(ut,oe),t(ut,mt),t(f,ie),t(f,M),t(M,Ot),t(M,at),t(M,St),t(M,O),t(M,bt),t(f,ae),t(f,z),t(z,Et),t(z,Vt),t(z,Nt),t(z,re),t(z,N),t(f,Wt),t(f,J),t(J,ht),t(J,ce),t(J,I),t(J,de),t(J,B),t(f,fe),t(f,P),t(P,qt),t(P,K),t(P,_t),t(P,pe),t(P,xt),t(f,ue),t(f,$),t($,Dt),t($,rt),t($,Ht),t($,me),t($,Mt),t(f,Xt),t(f,V),t(V,wt),t(V,be),t(V,It),t(V,he),t(V,Ct),t(f,_e),t(f,G),t(G,W),t(G,Yt),t(G,q),t(G,gt),t(G,A),t(f,xe),t(f,X),t(X,Y),t(X,y),t(X,Bt),t(X,Z),t(X,v),u(F,tt,$t),u(F,k,$t)},d(F){F&&m(n),F&&m(o),F&&m(a),F&&m(tt),F&&m(k)}}}function cl(c){let n,o,a,p,b;function d(_,f){return _[0]?rl:al}let h=d(c),C=h(c),x=c[0]&&Ie();return{c(){n=e("button"),C.c(),o=s(),x&&x.c(),a=Ye(),i(n,"class","btn btn-sm btn-secondary m-t-10")},m(_,f){u(_,n,f),C.m(n,null),u(_,o,f),x&&x.m(_,f),u(_,a,f),p||(b=Xe(n,"click",c[1]),p=!0)},p(_,[f]){h!==(h=d(_))&&(C.d(1),C=h(_),C&&(C.c(),C.m(n,null))),_[0]?x||(x=Ie(),x.c(),x.m(a.parentNode,a)):x&&(x.d(1),x=null)},i:He,o:He,d(_){_&&m(n),C.d(),_&&m(o),x&&x.d(_),_&&m(a),p=!1,b()}}}function dl(c,n,o){let a=!1;function p(){o(0,a=!a)}return[a,p]}class fl extends Ke{constructor(n){super(),Ve(this,n,dl,cl,We,{})}}function Be(c,n,o){const a=c.slice();return a[7]=n[o],a}function Ge(c,n,o){const a=c.slice();return a[7]=n[o],a}function Ue(c,n,o){const a=c.slice();return a[12]=n[o],a[14]=o,a}function je(c){let n;return{c(){n=e("p"),n.innerHTML="Requires admin <code>Authorization:TOKEN</code> header",i(n,"class","txt-hint txt-sm txt-right")},m(o,a){u(o,n,a)},d(o){o&&m(n)}}}function Qe(c){let n,o=c[12]+"",a,p=c[14]<c[4].length-1?", ":"",b;return{c(){n=e("code"),a=w(o),b=w(p)},m(d,h){u(d,n,h),t(n,a),u(d,b,h)},p(d,h){h&16&&o!==(o=d[12]+"")&&Ce(a,o),h&16&&p!==(p=d[14]<d[4].length-1?", ":"")&&Ce(b,p)},d(d){d&&m(n),d&&m(b)}}}function ze(c,n){let o,a=n[7].code+"",p,b,d,h;function C(){return n[6](n[7])}return{key:c,first:null,c(){o=e("div"),p=w(a),b=s(),i(o,"class","tab-item"),ke(o,"active",n[2]===n[7].code),this.first=o},m(x,_){u(x,o,_),t(o,p),t(o,b),d||(h=Xe(o,"click",C),d=!0)},p(x,_){n=x,_&36&&ke(o,"active",n[2]===n[7].code)},d(x){x&&m(o),d=!1,h()}}}function Je(c,n){let o,a,p,b;return a=new $e({props:{content:n[7].body}}),{key:c,first:null,c(){o=e("div"),Zt(a.$$.fragment),p=s(),i(o,"class","tab-item"),ke(o,"active",n[2]===n[7].code),this.first=o},m(d,h){u(d,o,h),te(a,o,null),t(o,p),b=!0},p(d,h){n=d,(!b||h&36)&&ke(o,"active",n[2]===n[7].code)},i(d){b||(Gt(a.$$.fragment,d),b=!0)},o(d){Ut(a.$$.fragment,d),b=!1},d(d){d&&m(o),ee(a)}}}function pl(c){var ye,Fe,Le,Ae,Te,Pe;let n,o,a=c[0].name+"",p,b,d,h,C,x,_,f=c[0].name+"",et,kt,jt,S,Qt,H,ct,R,lt,le,U,j,se,dt,vt=c[0].name+"",st,yt,ne,ft,pt,nt,E,zt,Ft,T,ot,Lt,Jt,At,Q,it,Tt,Kt,Pt,L,ut,Rt,oe,mt,ie,M,Ot,at,St,O,bt,ae,z,Et,Vt,Nt,re,N,Wt,J,ht,ce,I,de,B,fe,P,qt,K,_t,pe,xt,ue,$,Dt,rt,Ht,me,Mt,Xt,V,wt,be,It,he,Ct,_e,G,W,Yt,q,gt,A=[],xe=new Map,X,Y,y=[],Bt=new Map,Z;S=new il({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${c[3]}');

        ...

        // fetch a paginated records list
        const resultList = await pb.collection('${(ye=c[0])==null?void 0:ye.name}').getList(1, 50, {
            filter: 'created >= "2022-01-01 00:00:00" && someField1 != someField2',
        });

        // you can also fetch all records at once via getFullList
        const records = await pb.collection('${(Fe=c[0])==null?void 0:Fe.name}').getFullList({
            sort: '-created',
        });

        // or fetch only the first record that matches the specified filter
        const record = await pb.collection('${(Le=c[0])==null?void 0:Le.name}').getFirstListItem('someField="test"', {
            expand: 'relField1,relField2.subRelField',
        });
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${c[3]}');

        ...

        // fetch a paginated records list
        final resultList = await pb.collection('${(Ae=c[0])==null?void 0:Ae.name}').getList(
          page: 1,
          perPage: 50,
          filter: 'created >= "2022-01-01 00:00:00" && someField1 != someField2',
        );

        // you can also fetch all records at once via getFullList
        final records = await pb.collection('${(Te=c[0])==null?void 0:Te.name}').getFullList(
          sort: '-created',
        );

        // or fetch only the first record that matches the specified filter
        final record = await pb.collection('${(Pe=c[0])==null?void 0:Pe.name}').getFirstListItem(
          'someField="test"',
          expand: 'relField1,relField2.subRelField',
        );
    `}});let v=c[1]&&je();at=new $e({props:{content:`
                        // DESC by created and ASC by id
                        ?sort=-created,id
                    `}});let tt=c[4],k=[];for(let l=0;l<tt.length;l+=1)k[l]=Qe(Ue(c,tt,l));B=new $e({props:{content:`
                        ?filter=(id='abc' && created>'2022-01-01')
                    `}}),P=new fl({}),rt=new $e({props:{content:"?expand=relField1,relField2.subRelField"}});let F=c[5];const $t=l=>l[7].code;for(let l=0;l<F.length;l+=1){let r=Ge(c,F,l),g=$t(r);xe.set(g,A[l]=ze(g,r))}let we=c[5];const ve=l=>l[7].code;for(let l=0;l<we.length;l+=1){let r=Be(c,we,l),g=ve(r);Bt.set(g,y[l]=Je(g,r))}return{c(){n=e("h3"),o=w("List/Search ("),p=w(a),b=w(")"),d=s(),h=e("div"),C=e("p"),x=w("Fetch a paginated "),_=e("strong"),et=w(f),kt=w(" records list, supporting sorting and filtering."),jt=s(),Zt(S.$$.fragment),Qt=s(),H=e("h6"),H.textContent="API details",ct=s(),R=e("div"),lt=e("strong"),lt.textContent="GET",le=s(),U=e("div"),j=e("p"),se=w("/api/collections/"),dt=e("strong"),st=w(vt),yt=w("/records"),ne=s(),v&&v.c(),ft=s(),pt=e("div"),pt.textContent="Query parameters",nt=s(),E=e("table"),zt=e("thead"),zt.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr>`,Ft=s(),T=e("tbody"),ot=e("tr"),ot.innerHTML=`<td>page</td> 
            <td><span class="label">Number</span></td> 
            <td>The page (aka. offset) of the paginated list (default to 1).</td>`,Lt=s(),Jt=e("tr"),Jt.innerHTML=`<td>perPage</td> 
            <td><span class="label">Number</span></td> 
            <td>Specify the max returned records per page (default to 30).</td>`,At=s(),Q=e("tr"),it=e("td"),it.textContent="sort",Tt=s(),Kt=e("td"),Kt.innerHTML='<span class="label">String</span>',Pt=s(),L=e("td"),ut=w("Specify the records order attribute(s). "),Rt=e("br"),oe=w(`
                Add `),mt=e("code"),mt.textContent="-",ie=w(" / "),M=e("code"),M.textContent="+",Ot=w(` (default) in front of the attribute for DESC / ASC order.
                Ex.:
                `),Zt(at.$$.fragment),St=s(),O=e("p"),bt=e("strong"),bt.textContent="Supported record sort fields:",ae=s(),z=e("br"),Et=s(),Vt=e("code"),Vt.textContent="@random",Nt=w(`,
                    `);for(let l=0;l<k.length;l+=1)k[l].c();re=s(),N=e("tr"),Wt=e("td"),Wt.textContent="filter",J=s(),ht=e("td"),ht.innerHTML='<span class="label">String</span>',ce=s(),I=e("td"),de=w(`Filter the returned records. Ex.:
                `),Zt(B.$$.fragment),fe=s(),Zt(P.$$.fragment),qt=s(),K=e("tr"),_t=e("td"),_t.textContent="expand",pe=s(),xt=e("td"),xt.innerHTML='<span class="label">String</span>',ue=s(),$=e("td"),Dt=w(`Auto expand record relations. Ex.:
                `),Zt(rt.$$.fragment),Ht=w(`
                Supports up to 6-levels depth nested relations expansion. `),me=e("br"),Mt=w(`
                The expanded relations will be appended to each individual record under the
                `),Xt=e("code"),Xt.textContent="expand",V=w(" property (eg. "),wt=e("code"),wt.textContent='"expand": {"relField1": {...}, ...}',be=w(`).
                `),It=e("br"),he=w(`
                Only the relations to which the request user has permissions to `),Ct=e("strong"),Ct.textContent="view",_e=w(" will be expanded."),G=s(),W=e("div"),W.textContent="Responses",Yt=s(),q=e("div"),gt=e("div");for(let l=0;l<A.length;l+=1)A[l].c();X=s(),Y=e("div");for(let l=0;l<y.length;l+=1)y[l].c();i(n,"class","m-b-sm"),i(h,"class","content txt-lg m-b-sm"),i(H,"class","m-b-xs"),i(lt,"class","label label-primary"),i(U,"class","content"),i(R,"class","alert alert-info"),i(pt,"class","section-title"),i(E,"class","table-compact table-border m-b-base"),i(W,"class","section-title"),i(gt,"class","tabs-header compact left"),i(Y,"class","tabs-content"),i(q,"class","tabs")},m(l,r){u(l,n,r),t(n,o),t(n,p),t(n,b),u(l,d,r),u(l,h,r),t(h,C),t(C,x),t(C,_),t(_,et),t(C,kt),u(l,jt,r),te(S,l,r),u(l,Qt,r),u(l,H,r),u(l,ct,r),u(l,R,r),t(R,lt),t(R,le),t(R,U),t(U,j),t(j,se),t(j,dt),t(dt,st),t(j,yt),t(R,ne),v&&v.m(R,null),u(l,ft,r),u(l,pt,r),u(l,nt,r),u(l,E,r),t(E,zt),t(E,Ft),t(E,T),t(T,ot),t(T,Lt),t(T,Jt),t(T,At),t(T,Q),t(Q,it),t(Q,Tt),t(Q,Kt),t(Q,Pt),t(Q,L),t(L,ut),t(L,Rt),t(L,oe),t(L,mt),t(L,ie),t(L,M),t(L,Ot),te(at,L,null),t(L,St),t(L,O),t(O,bt),t(O,ae),t(O,z),t(O,Et),t(O,Vt),t(O,Nt);for(let g=0;g<k.length;g+=1)k[g]&&k[g].m(O,null);t(T,re),t(T,N),t(N,Wt),t(N,J),t(N,ht),t(N,ce),t(N,I),t(I,de),te(B,I,null),t(I,fe),te(P,I,null),t(T,qt),t(T,K),t(K,_t),t(K,pe),t(K,xt),t(K,ue),t(K,$),t($,Dt),te(rt,$,null),t($,Ht),t($,me),t($,Mt),t($,Xt),t($,V),t($,wt),t($,be),t($,It),t($,he),t($,Ct),t($,_e),u(l,G,r),u(l,W,r),u(l,Yt,r),u(l,q,r),t(q,gt);for(let g=0;g<A.length;g+=1)A[g]&&A[g].m(gt,null);t(q,X),t(q,Y);for(let g=0;g<y.length;g+=1)y[g]&&y[g].m(Y,null);Z=!0},p(l,[r]){var Re,Oe,Se,Ee,Ne,qe;(!Z||r&1)&&a!==(a=l[0].name+"")&&Ce(p,a),(!Z||r&1)&&f!==(f=l[0].name+"")&&Ce(et,f);const g={};if(r&9&&(g.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${l[3]}');

        ...

        // fetch a paginated records list
        const resultList = await pb.collection('${(Re=l[0])==null?void 0:Re.name}').getList(1, 50, {
            filter: 'created >= "2022-01-01 00:00:00" && someField1 != someField2',
        });

        // you can also fetch all records at once via getFullList
        const records = await pb.collection('${(Oe=l[0])==null?void 0:Oe.name}').getFullList({
            sort: '-created',
        });

        // or fetch only the first record that matches the specified filter
        const record = await pb.collection('${(Se=l[0])==null?void 0:Se.name}').getFirstListItem('someField="test"', {
            expand: 'relField1,relField2.subRelField',
        });
    `),r&9&&(g.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${l[3]}');

        ...

        // fetch a paginated records list
        final resultList = await pb.collection('${(Ee=l[0])==null?void 0:Ee.name}').getList(
          page: 1,
          perPage: 50,
          filter: 'created >= "2022-01-01 00:00:00" && someField1 != someField2',
        );

        // you can also fetch all records at once via getFullList
        final records = await pb.collection('${(Ne=l[0])==null?void 0:Ne.name}').getFullList(
          sort: '-created',
        );

        // or fetch only the first record that matches the specified filter
        final record = await pb.collection('${(qe=l[0])==null?void 0:qe.name}').getFirstListItem(
          'someField="test"',
          expand: 'relField1,relField2.subRelField',
        );
    `),S.$set(g),(!Z||r&1)&&vt!==(vt=l[0].name+"")&&Ce(st,vt),l[1]?v||(v=je(),v.c(),v.m(R,null)):v&&(v.d(1),v=null),r&16){tt=l[4];let D;for(D=0;D<tt.length;D+=1){const De=Ue(l,tt,D);k[D]?k[D].p(De,r):(k[D]=Qe(De),k[D].c(),k[D].m(O,null))}for(;D<k.length;D+=1)k[D].d(1);k.length=tt.length}r&36&&(F=l[5],A=Me(A,r,$t,1,l,F,xe,gt,Ze,ze,null,Ge)),r&36&&(we=l[5],tl(),y=Me(y,r,ve,1,l,we,Bt,Y,el,Je,null,Be),ll())},i(l){if(!Z){Gt(S.$$.fragment,l),Gt(at.$$.fragment,l),Gt(B.$$.fragment,l),Gt(P.$$.fragment,l),Gt(rt.$$.fragment,l);for(let r=0;r<we.length;r+=1)Gt(y[r]);Z=!0}},o(l){Ut(S.$$.fragment,l),Ut(at.$$.fragment,l),Ut(B.$$.fragment,l),Ut(P.$$.fragment,l),Ut(rt.$$.fragment,l);for(let r=0;r<y.length;r+=1)Ut(y[r]);Z=!1},d(l){l&&m(n),l&&m(d),l&&m(h),l&&m(jt),ee(S,l),l&&m(Qt),l&&m(H),l&&m(ct),l&&m(R),v&&v.d(),l&&m(ft),l&&m(pt),l&&m(nt),l&&m(E),ee(at),sl(k,l),ee(B),ee(P),ee(rt),l&&m(G),l&&m(W),l&&m(Yt),l&&m(q);for(let r=0;r<A.length;r+=1)A[r].d();for(let r=0;r<y.length;r+=1)y[r].d()}}}function ul(c,n,o){let a,p,b,{collection:d=new nl}=n,h=200,C=[];const x=_=>o(2,h=_.code);return c.$$set=_=>{"collection"in _&&o(0,d=_.collection)},c.$$.update=()=>{c.$$.dirty&1&&o(4,a=ge.getAllCollectionIdentifiers(d)),c.$$.dirty&1&&o(1,p=(d==null?void 0:d.listRule)===null),c.$$.dirty&3&&d!=null&&d.id&&(C.push({code:200,body:JSON.stringify({page:1,perPage:30,totalPages:1,totalItems:2,items:[ge.dummyCollectionRecord(d),ge.dummyCollectionRecord(d)]},null,2)}),C.push({code:400,body:`
                {
                  "code": 400,
                  "message": "Something went wrong while processing your request. Invalid filter.",
                  "data": {}
                }
            `}),p&&C.push({code:403,body:`
                    {
                      "code": 403,
                      "message": "Only admins can access this action.",
                      "data": {}
                    }
                `}))},o(3,b=ge.getApiExampleUrl(ol.baseUrl)),[d,p,h,b,a,C,x]}class hl extends Ke{constructor(n){super(),Ve(this,n,ul,pl,We,{collection:0})}}export{hl as default};
