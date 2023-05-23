import{S as tt,i as lt,s as st,M as et,e as o,w as b,b as u,c as W,f as _,g as r,h as l,m as X,x as ve,N as Ge,P as nt,k as ot,Q as it,n as at,t as U,a as j,o as d,d as Y,T as rt,C as Je,p as dt,r as Z,u as ct}from"./index-a65ca895.js";import{S as ft}from"./SdkTabs-ad912c8f.js";import{F as pt}from"./FieldsQueryParam-ba250473.js";function Ke(i,s,n){const a=i.slice();return a[6]=s[n],a}function We(i,s,n){const a=i.slice();return a[6]=s[n],a}function Xe(i){let s;return{c(){s=o("p"),s.innerHTML="Requires admin <code>Authorization:TOKEN</code> header",_(s,"class","txt-hint txt-sm txt-right")},m(n,a){r(n,s,a)},d(n){n&&d(s)}}}function Ye(i,s){let n,a=s[6].code+"",w,c,f,m;function F(){return s[5](s[6])}return{key:i,first:null,c(){n=o("button"),w=b(a),c=u(),_(n,"class","tab-item"),Z(n,"active",s[2]===s[6].code),this.first=n},m(h,g){r(h,n,g),l(n,w),l(n,c),f||(m=ct(n,"click",F),f=!0)},p(h,g){s=h,g&20&&Z(n,"active",s[2]===s[6].code)},d(h){h&&d(n),f=!1,m()}}}function Ze(i,s){let n,a,w,c;return a=new et({props:{content:s[6].body}}),{key:i,first:null,c(){n=o("div"),W(a.$$.fragment),w=u(),_(n,"class","tab-item"),Z(n,"active",s[2]===s[6].code),this.first=n},m(f,m){r(f,n,m),X(a,n,null),l(n,w),c=!0},p(f,m){s=f,(!c||m&20)&&Z(n,"active",s[2]===s[6].code)},i(f){c||(U(a.$$.fragment,f),c=!0)},o(f){j(a.$$.fragment,f),c=!1},d(f){f&&d(n),Y(a)}}}function ut(i){var Ue,je;let s,n,a=i[0].name+"",w,c,f,m,F,h,g,V=i[0].name+"",ee,$e,te,R,le,M,se,y,z,we,G,E,ye,ne,J=i[0].name+"",oe,Ce,ie,Fe,ae,x,re,A,de,I,ce,O,fe,ge,q,P,pe,Re,ue,Oe,k,Pe,S,Te,De,Ee,me,Se,be,Be,Me,xe,_e,Ae,Ie,B,ke,H,he,T,L,C=[],qe=new Map,He,N,v=[],Le=new Map,D;R=new ft({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${i[3]}');

        ...

        const record = await pb.collection('${(Ue=i[0])==null?void 0:Ue.name}').getOne('RECORD_ID', {
            expand: 'relField1,relField2.subRelField',
        });
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${i[3]}');

        ...

        final record = await pb.collection('${(je=i[0])==null?void 0:je.name}').getOne('RECORD_ID',
          expand: 'relField1,relField2.subRelField',
        );
    `}});let $=i[1]&&Xe();S=new et({props:{content:"?expand=relField1,relField2.subRelField"}}),B=new pt({});let K=i[4];const Ne=e=>e[6].code;for(let e=0;e<K.length;e+=1){let t=We(i,K,e),p=Ne(t);qe.set(p,C[e]=Ye(p,t))}let Q=i[4];const Qe=e=>e[6].code;for(let e=0;e<Q.length;e+=1){let t=Ke(i,Q,e),p=Qe(t);Le.set(p,v[e]=Ze(p,t))}return{c(){s=o("h3"),n=b("View ("),w=b(a),c=b(")"),f=u(),m=o("div"),F=o("p"),h=b("Fetch a single "),g=o("strong"),ee=b(V),$e=b(" record."),te=u(),W(R.$$.fragment),le=u(),M=o("h6"),M.textContent="API details",se=u(),y=o("div"),z=o("strong"),z.textContent="GET",we=u(),G=o("div"),E=o("p"),ye=b("/api/collections/"),ne=o("strong"),oe=b(J),Ce=b("/records/"),ie=o("strong"),ie.textContent=":id",Fe=u(),$&&$.c(),ae=u(),x=o("div"),x.textContent="Path Parameters",re=u(),A=o("table"),A.innerHTML=`<thead><tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr></thead> 
    <tbody><tr><td>id</td> 
            <td><span class="label">String</span></td> 
            <td>ID of the record to view.</td></tr></tbody>`,de=u(),I=o("div"),I.textContent="Query parameters",ce=u(),O=o("table"),fe=o("thead"),fe.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr>`,ge=u(),q=o("tbody"),P=o("tr"),pe=o("td"),pe.textContent="expand",Re=u(),ue=o("td"),ue.innerHTML='<span class="label">String</span>',Oe=u(),k=o("td"),Pe=b(`Auto expand record relations. Ex.:
                `),W(S.$$.fragment),Te=b(`
                Supports up to 6-levels depth nested relations expansion. `),De=o("br"),Ee=b(`
                The expanded relations will be appended to the record under the
                `),me=o("code"),me.textContent="expand",Se=b(" property (eg. "),be=o("code"),be.textContent='"expand": {"relField1": {...}, ...}',Be=b(`).
                `),Me=o("br"),xe=b(`
                Only the relations to which the request user has permissions to `),_e=o("strong"),_e.textContent="view",Ae=b(" will be expanded."),Ie=u(),W(B.$$.fragment),ke=u(),H=o("div"),H.textContent="Responses",he=u(),T=o("div"),L=o("div");for(let e=0;e<C.length;e+=1)C[e].c();He=u(),N=o("div");for(let e=0;e<v.length;e+=1)v[e].c();_(s,"class","m-b-sm"),_(m,"class","content txt-lg m-b-sm"),_(M,"class","m-b-xs"),_(z,"class","label label-primary"),_(G,"class","content"),_(y,"class","alert alert-info"),_(x,"class","section-title"),_(A,"class","table-compact table-border m-b-base"),_(I,"class","section-title"),_(O,"class","table-compact table-border m-b-base"),_(H,"class","section-title"),_(L,"class","tabs-header compact left"),_(N,"class","tabs-content"),_(T,"class","tabs")},m(e,t){r(e,s,t),l(s,n),l(s,w),l(s,c),r(e,f,t),r(e,m,t),l(m,F),l(F,h),l(F,g),l(g,ee),l(F,$e),r(e,te,t),X(R,e,t),r(e,le,t),r(e,M,t),r(e,se,t),r(e,y,t),l(y,z),l(y,we),l(y,G),l(G,E),l(E,ye),l(E,ne),l(ne,oe),l(E,Ce),l(E,ie),l(y,Fe),$&&$.m(y,null),r(e,ae,t),r(e,x,t),r(e,re,t),r(e,A,t),r(e,de,t),r(e,I,t),r(e,ce,t),r(e,O,t),l(O,fe),l(O,ge),l(O,q),l(q,P),l(P,pe),l(P,Re),l(P,ue),l(P,Oe),l(P,k),l(k,Pe),X(S,k,null),l(k,Te),l(k,De),l(k,Ee),l(k,me),l(k,Se),l(k,be),l(k,Be),l(k,Me),l(k,xe),l(k,_e),l(k,Ae),l(q,Ie),X(B,q,null),r(e,ke,t),r(e,H,t),r(e,he,t),r(e,T,t),l(T,L);for(let p=0;p<C.length;p+=1)C[p]&&C[p].m(L,null);l(T,He),l(T,N);for(let p=0;p<v.length;p+=1)v[p]&&v[p].m(N,null);D=!0},p(e,[t]){var Ve,ze;(!D||t&1)&&a!==(a=e[0].name+"")&&ve(w,a),(!D||t&1)&&V!==(V=e[0].name+"")&&ve(ee,V);const p={};t&9&&(p.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        const record = await pb.collection('${(Ve=e[0])==null?void 0:Ve.name}').getOne('RECORD_ID', {
            expand: 'relField1,relField2.subRelField',
        });
    `),t&9&&(p.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        final record = await pb.collection('${(ze=e[0])==null?void 0:ze.name}').getOne('RECORD_ID',
          expand: 'relField1,relField2.subRelField',
        );
    `),R.$set(p),(!D||t&1)&&J!==(J=e[0].name+"")&&ve(oe,J),e[1]?$||($=Xe(),$.c(),$.m(y,null)):$&&($.d(1),$=null),t&20&&(K=e[4],C=Ge(C,t,Ne,1,e,K,qe,L,nt,Ye,null,We)),t&20&&(Q=e[4],ot(),v=Ge(v,t,Qe,1,e,Q,Le,N,it,Ze,null,Ke),at())},i(e){if(!D){U(R.$$.fragment,e),U(S.$$.fragment,e),U(B.$$.fragment,e);for(let t=0;t<Q.length;t+=1)U(v[t]);D=!0}},o(e){j(R.$$.fragment,e),j(S.$$.fragment,e),j(B.$$.fragment,e);for(let t=0;t<v.length;t+=1)j(v[t]);D=!1},d(e){e&&d(s),e&&d(f),e&&d(m),e&&d(te),Y(R,e),e&&d(le),e&&d(M),e&&d(se),e&&d(y),$&&$.d(),e&&d(ae),e&&d(x),e&&d(re),e&&d(A),e&&d(de),e&&d(I),e&&d(ce),e&&d(O),Y(S),Y(B),e&&d(ke),e&&d(H),e&&d(he),e&&d(T);for(let t=0;t<C.length;t+=1)C[t].d();for(let t=0;t<v.length;t+=1)v[t].d()}}}function mt(i,s,n){let a,w,{collection:c=new rt}=s,f=200,m=[];const F=h=>n(2,f=h.code);return i.$$set=h=>{"collection"in h&&n(0,c=h.collection)},i.$$.update=()=>{i.$$.dirty&1&&n(1,a=(c==null?void 0:c.viewRule)===null),i.$$.dirty&3&&c!=null&&c.id&&(m.push({code:200,body:JSON.stringify(Je.dummyCollectionRecord(c),null,2)}),a&&m.push({code:403,body:`
                    {
                      "code": 403,
                      "message": "Only admins can access this action.",
                      "data": {}
                    }
                `}),m.push({code:404,body:`
                {
                  "code": 404,
                  "message": "The requested resource wasn't found.",
                  "data": {}
                }
            `}))},n(3,w=Je.getApiExampleUrl(dt.baseUrl)),[c,a,f,w,m,F]}class ht extends tt{constructor(s){super(),lt(this,s,mt,ut,st,{collection:0})}}export{ht as default};
