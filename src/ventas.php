<?php include_once "includes/header.php";
include "../conexion.php";
$id_user = $_SESSION['idUser'];
$permiso = "nueva_venta";
$sql = mysqli_query($conexion, "SELECT p.*, d.* FROM permisos p INNER JOIN detalle_permisos d ON p.id = d.id_permiso WHERE d.id_usuario = $id_user AND p.nombre = '$permiso'");
$existe = mysqli_fetch_all($sql);
if (empty($existe) && $id_user != 1) {
    header("Location: permisos.php");
}

?>
<div class="row">
    <div class="col-lg-12">
        <div class="form-group">
            <h4 class="text-center">Datos del Cliente</h4>
        </div>
        <div class="card">
            <div class="card-body">
                <form method="post">
                    <div class="row">
                        <div class="col-lg-3">
                            <div>
                                <input type="hidden" id="idcliente" value="1" name="idcliente" required>
                                <label>Nombre Cliente</label>
                                <input type="text" name="nom_cliente" id="nom_cliente" class="form-control" placeholder="Ingresá el nombre" required>
                            </div>
                        </div>
                        <div class="col-lg-3">
                            <div class="form-group">
                                <label>Teléfono</label>
                                <input type="number" name="tel_cliente" id="tel_cliente" class="form-control" disabled required>
                            </div>
                        </div>
                        <div class="col-lg-3">
                            <div class="form-group">
                                <label>Dirreción</label>
                                <input type="text" name="dir_cliente" id="dir_cliente" class="form-control" disabled required>
                            </div>
                        </div>
                        <div class="col-lg-3">
                            <div class="form-group">
                                <label>Obra Social</label>
                                <input type="text" name="obrasocial" id="obrasocial" class="form-control" disabled required>
                            </div>
                        </div>
                    </div>
                </form>
            </div>
        </div>
        <div class="card">
            <div class="card-header bg-primary text-white text-center">
                Datos Venta
            </div>
            <div class="card-body">
                <div class="row">
                    <div class="col-lg-6">
                        <div class="form-group">
                            Vendedor <i class="fas fa-user"></i><p style="font-size: 16px; text-transform: uppercase; color: red;"><?php echo $_SESSION['nombre']; ?></p>
                            <label> <b>Graduacion Lejos &nbsp&nbsp&nbsp&nbsp&nbsp&nbsp&nbsp&nbsp
                                &nbsp&nbsp&nbsp&nbsp&nbsp&nbsp&nbsp&nbsp&nbsp&nbsp&nbsp&nbsp&nbsp
                                &nbsp&nbsp&nbsp&nbsp&nbsp&nbsp&nbsp&nbsp&nbsp&nbsp&nbsp&nbsp&nbsp&nbsp
                                &nbsp&nbsp&nbsp Graduacion Cerca</b> </label>
                            <p style="font-size: 16px; text-transform: uppercase; color: black;">
                            <table class="table table-borderless" id="tablaGracuadiones">
                            <form id="graduaciones">


                                
                                
                                <div class="input-group">    
                                <span class="input-group-text">Ojo DL</span>
                                <input type="text" id="ojoDl1" name="ojoDl1" class="form-control" />
                                <input type="text" id="ojoDl2" name="ojoDl2" class="form-control" />
                                <input type="text" id="ojoDl3" name="ojoDl3" class="form-control" />
                                <span class="input-group-text">Ojo DC  </span>
                                <input type="text" name="ojoD1" id="ojoD1" class="form-control" />
                                <input type="text" name="ojoD2" id="ojoD2" class="form-control" />
                                <input type="text" name="ojoD3" id="ojoD3" class="form-control" />
                                </div><br>
                                <div class="input-group">
                                <span class="input-group-text">Ojo Iz L</span>
                                <input type="text" id="ojoIl1" name="ojoIl1" class="form-control" />
                                <input type="text" id="ojoIl2" name="ojoIl2" class="form-control" />
                                <input type="text" id="ojoIl3" name="ojoIl3" class="form-control" />
                                <span class="input-group-text">Ojo Iz  C</span>
                                <input type="text" id="ojoI1" name="ojoI1" class="form-control" />
                                <input type="text" id="ojoI2" name="ojoI2" class="form-control" />
                                <input type="text" id="ojoI3" name="ojoI3" class="form-control" />
                                </div><br>
                                <div class="input-group">
                                <span class="input-group-text">ADD:</span>
                                <input type="text" id="add" name="add" class="form-control" />
                                <span class="input-group-text">Obse:</span>
                                <input type="text" id="obs" name="obs" class="form-control" />
                                </div><br>
                                
                                <tr>
                                    <td><input class="btn btn-primary" name="grad" id="grad" type="button" value="Agregar Graduaciones"></td>
                                    <td><form type="post"id="borrar_grad"><input class="btn btn-danger" type="hidden" value="Borrar Graduaciones"id="borrar_grad"></form></td>
                                    
                                </tr>
                              
                                
                            
                            </form>
                        </table>   
                        
                        </p>
                        </div>
                    </div>
                    
                    <div class="col-lg-6">
                        <div class="card">
                            <div class="card-header">
                                Buscar Producto
                            </div>
                            <div class="card-body">
                                <div class="form-group">
                                    <input id="producto" class="form-control" type="text" name="producto" placeholder="Ingresá el código o nombre">
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
        
            
        </div>
        <div id="okgrad"></div>
        <div class="table-responsive">
            <table class="table table-hover" id="tblDetalle">
                <thead class="thead-dark">
                    <tr>
                        <th>Id</th>
                        <th>Descripción</th>
                        <th>Cantidad</th>
                        <th>Precio</th>
                        <th>Precio Total</th>
                        <th>Accion</th>
                        <th></th>
                        <th></th>
                    </tr>
                </thead>
                <tbody id="detalle_venta">
                </tbody>
                <tfoot>
                    <tr class="font-weight-bold">
                        <td colspan=3>Total a Pagar: </td>
                        <td></td>
                    </tr>    
                </tfoot>
            </table>
            <table class="table table-hover">
                <form method="POST" id="form_descuento">
                    <div id="div_descuento">
                    <tr class="font-weight-bold">
                        <td colspan=3>Abona: </td>
                        <td colspan=3><input type="number" size="3" id="abona"> </td>
                    </tr>
                    <tr class="font-weight-bold">
                        <td colspan=3>Descuento: </td>
                        <td colspan=3><select id="porc" name="porc">
                                        <option value="1">Sin descuento</option>
                                        <option value="0.95">5%</option>
                                        <option value="0.9">10%</option>
                                        <option value="0.85">15%</option>
                                        <option value="0.80">20%</option>
                                        <option value="0.75">25%</option>
                                        <option value="0.70">30%</option>
                                        <option value="0.65">35%</option>
                                        <option value="0.60">40%</option>
                                        <option value="0.55">45%</option>
                                        <option value="0.50">50%</option>
                                        <option value="0.45">55%</option>
                                        <option value="0.40">60%</option>
                                        </select></td>
                    </tr>
                    <tr class="font-weight-bold">
                        <td colspan=3>Obra Social: </td>
                        <td colspan=3><input type="number" size="3" id="obra_social">
                        <input type="number" size="3" id="total" hidden disabled></td>
                    </tr>
                    <tr class="font-weight-bold">
                        <td colspan=3>Resta: </td>
                        <td colspan=3><input type="number" size="3" id="resto" disabled></td>
                    </tr>
                </div>
                </form>
                    </tr>
                    
                </tfoot>
            </table>
        </div>
    </div>
    <div class="col-md-6"><br>
        <form method="POST" id="metodo_pago">
            <h4 class="text-left">Método de Pago</h4>
            <label for="male">Efectivo </label>
            <input type="radio" id="1" name="pago" value="1" checked>
            <label for="female">&nbsp&nbsp Tarjeta de Crédito </label>
            <input type="radio" id="2" name="pago" value="2">
            <label for="other">&nbsp&nbsp Tarjeta de Débito </label>
            <input type="radio" id="3" name="pago" value="3">
            <label for="other">&nbsp&nbsp Transferencia </label>
            <input type="radio" id="4" name="pago" value="4">        
        </form>
    </div>
    <div class="col-md-6"><br>
        <a href="#" class="btn btn-primary" id="btn_generar"><i class="fas fa-save"></i> Generar Venta</a>
        <input type="button" class="btn btn-primary" value="Simular Venta" name="btn_parcial" id="btn_parcial" ></input>
    </div>
</div>
<?php include_once "includes/footer.php"; ?>