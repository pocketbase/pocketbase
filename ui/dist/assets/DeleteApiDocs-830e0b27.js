import{S as Ce,i as Re,s as Pe,e as c,w as $,b as h,c as $e,f as m,g as f,h as n,m as ge,x,O as _e,P as Ee,k as Oe,Q as Te,n as Be,t as ee,a as te,o as p,d as we,C as Ie,p as Ae,r as N,u as Me,N as Se}from"./index-cc2b3041.js";import{S as qe}from"./SdkTabs-cf23198f.js";function ke(o,l,s){const a=o.slice();return a[6]=l[s],a}function he(o,l,s){const a=o.slice();return a[6]=l[s],a}function ve(o){let l;return{c(){l=c("p"),l.innerHTML="Requires admin <code>Authorization:TOKEN</code> header",m(l,"class","txt-hint txt-sm txt-right")},m(s,a){f(s,l,a)},d(s){s&&p(l)}}}function ye(o,l){let s,a=l[6].code+"",v,i,r,u;function g(){return l[5](l[6])}return{key:o,first:null,c(){s=c("button"),v=$(a),i=h(),m(s,"class","tab-item"),N(s,"active",l[2]===l[6].code),this.first=s},m(b,w){f(b,s,w),n(s,v),n(s,i),r||(u=Me(s,"click",g),r=!0)},p(b,w){l=b,w&20&&N(s,"active",l[2]===l[6].code)},d(b){b&&p(s),r=!1,u()}}}function De(o,l){let s,a,v,i;return a=new Se({props:{content:l[6].body}}),{key:o,first:null,c(){s=c("div"),$e(a.$$.fragment),v=h(),m(s,"class","tab-item"),N(s,"active",l[2]===l[6].code),this.first=s},m(r,u){f(r,s,u),ge(a,s,null),n(s,v),i=!0},p(r,u){l=r,(!i||u&20)&&N(s,"active",l[2]===l[6].code)},i(r){i||(ee(a.$$.fragment,r),i=!0)},o(r){te(a.$$.fragment,r),i=!1},d(r){r&&p(s),we(a)}}}function He(o){var pe,ue;let l,s,a=o[0].name+"",v,i,r,u,g,b,w,q=o[0].name+"",z,le,F,C,K,O,Q,y,H,se,L,E,oe,G,U=o[0].name+"",J,ae,V,ne,W,T,X,B,Y,I,Z,R,A,D=[],ie=new Map,re,M,_=[],ce=new Map,P;C=new qe({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${o[3]}');

        ...

        await pb.collection('${(pe=o[0])==null?void 0:pe.name}').delete('RECORD_ID');
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${o[3]}');

        ...

        await pb.collection('${(ue=o[0])==null?void 0:ue.name}').delete('RECORD_ID');
    `}});let k=o[1]&&ve(),j=o[4];const de=e=>e[6].code;for(let e=0;e<j.length;e+=1){let t=he(o,j,e),d=de(t);ie.set(d,D[e]=ye(d,t))}let S=o[4];const fe=e=>e[6].code;for(let e=0;e<S.length;e+=1){let t=ke(o,S,e),d=fe(t);ce.set(d,_[e]=De(d,t))}return{c(){l=c("h3"),s=$("Delete ("),v=$(a),i=$(")"),r=h(),u=c("div"),g=c("p"),b=$("Delete a single "),w=c("strong"),z=$(q),le=$(" record."),F=h(),$e(C.$$.fragment),K=h(),O=c("h6"),O.textContent="API details",Q=h(),y=c("div"),H=c("strong"),H.textContent="DELETE",se=h(),L=c("div"),E=c("p"),oe=$("/api/collections/"),G=c("strong"),J=$(U),ae=$("/records/"),V=c("strong"),V.textContent=":id",ne=h(),k&&k.c(),W=h(),T=c("div"),T.textContent="Path parameters",X=h(),B=c("table"),B.innerHTML=`<thead><tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr></thead> 
    <tbody><tr><td>id</td> 
            <td><span class="label">String</span></td> 
            <td>ID of the record to delete.</td></tr></tbody>`,Y=h(),I=c("div"),I.textContent="Responses",Z=h(),R=c("div"),A=c("div");for(let e=0;e<D.length;e+=1)D[e].c();re=h(),M=c("div");for(let e=0;e<_.length;e+=1)_[e].c();m(l,"class","m-b-sm"),m(u,"class","content txt-lg m-b-sm"),m(O,"class","m-b-xs"),m(H,"class","label label-primary"),m(L,"class","content"),m(y,"class","alert alert-danger"),m(T,"class","section-title"),m(B,"class","table-compact table-border m-b-base"),m(I,"class","section-title"),m(A,"class","tabs-header compact left"),m(M,"class","tabs-content"),m(R,"class","tabs")},m(e,t){f(e,l,t),n(l,s),n(l,v),n(l,i),f(e,r,t),f(e,u,t),n(u,g),n(g,b),n(g,w),n(w,z),n(g,le),f(e,F,t),ge(C,e,t),f(e,K,t),f(e,O,t),f(e,Q,t),f(e,y,t),n(y,H),n(y,se),n(y,L),n(L,E),n(E,oe),n(E,G),n(G,J),n(E,ae),n(E,V),n(y,ne),k&&k.m(y,null),f(e,W,t),f(e,T,t),f(e,X,t),f(e,B,t),f(e,Y,t),f(e,I,t),f(e,Z,t),f(e,R,t),n(R,A);for(let d=0;d<D.length;d+=1)D[d]&&D[d].m(A,null);n(R,re),n(R,M);for(let d=0;d<_.length;d+=1)_[d]&&_[d].m(M,null);P=!0},p(e,[t]){var me,be;(!P||t&1)&&a!==(a=e[0].name+"")&&x(v,a),(!P||t&1)&&q!==(q=e[0].name+"")&&x(z,q);const d={};t&9&&(d.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(me=e[0])==null?void 0:me.name}').delete('RECORD_ID');
    `),t&9&&(d.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${(be=e[0])==null?void 0:be.name}').delete('RECORD_ID');
    `),C.$set(d),(!P||t&1)&&U!==(U=e[0].name+"")&&x(J,U),e[1]?k||(k=ve(),k.c(),k.m(y,null)):k&&(k.d(1),k=null),t&20&&(j=e[4],D=_e(D,t,de,1,e,j,ie,A,Ee,ye,null,he)),t&20&&(S=e[4],Oe(),_=_e(_,t,fe,1,e,S,ce,M,Te,De,null,ke),Be())},i(e){if(!P){ee(C.$$.fragment,e);for(let t=0;t<S.length;t+=1)ee(_[t]);P=!0}},o(e){te(C.$$.fragment,e);for(let t=0;t<_.length;t+=1)te(_[t]);P=!1},d(e){e&&p(l),e&&p(r),e&&p(u),e&&p(F),we(C,e),e&&p(K),e&&p(O),e&&p(Q),e&&p(y),k&&k.d(),e&&p(W),e&&p(T),e&&p(X),e&&p(B),e&&p(Y),e&&p(I),e&&p(Z),e&&p(R);for(let t=0;t<D.length;t+=1)D[t].d();for(let t=0;t<_.length;t+=1)_[t].d()}}}function Le(o,l,s){let a,v,{collection:i}=l,r=204,u=[];const g=b=>s(2,r=b.code);return o.$$set=b=>{"collection"in b&&s(0,i=b.collection)},o.$$.update=()=>{o.$$.dirty&1&&s(1,a=(i==null?void 0:i.deleteRule)===null),o.$$.dirty&3&&i!=null&&i.id&&(u.push({code:204,body:`
                null
            `}),u.push({code:400,body:`
                {
                  "code": 400,
                  "message": "Failed to delete record. Make sure that the record is not part of a required relation reference.",
                  "data": {}
                }
            `}),a&&u.push({code:403,body:`
                    {
                      "code": 403,
                      "message": "Only admins can access this action.",
                      "data": {}
                    }
                `}),u.push({code:404,body:`
                {
                  "code": 404,
                  "message": "The requested resource wasn't found.",
                  "data": {}
                }
            `}))},s(3,v=Ie.getApiExampleUrl(Ae.baseUrl)),[i,a,r,v,u,g]}class Ne extends Ce{constructor(l){super(),Re(this,l,Le,He,Pe,{collection:0})}}export{Ne as default};
